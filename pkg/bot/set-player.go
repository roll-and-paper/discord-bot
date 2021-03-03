package bot

import (
	"context"
	"github.com/bwmarrin/discordgo"
	gp "github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
	"github.com/thoas/go-funk"
	"strings"
)

func setPlayer(s *Services, msg *discordgo.Message, sess *discordgo.Session, state *gp.ServerState, args []string) error {
	if len(args) < 2 {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.player.missing", state)
		return err
	}
	guild, err := sess.Guild(msg.GuildID)
	if err != nil {
		return err
	}

	playerRole := s.GetRole(guild, state.Roles.Players)
	if playerRole == nil {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.player.role-missing", state)
		return err
	}

	playersChannel, err := s.GetChannel(sess, state.Channels.Players)
	if err != nil {
		return err
	} else if playersChannel == nil {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.player.channel-missing", state)
		return err
	}

	maybePlayerId := discord.GetUserIdFromMention(args[0])
	if maybePlayerId == "" {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.not-a-person", map[string]interface{}{"Name": args[0]})
		return err
	}

	maybePlayer, err := sess.GuildMember(msg.GuildID, maybePlayerId)
	if err != nil || maybePlayer == nil {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.unknown-person", map[string]interface{}{"Name": maybePlayerId})
		return err
	}

	if err := sess.GuildMemberRoleAdd(guild.ID, maybePlayer.User.ID, playerRole.ID); err != nil {
		return err
	}

	characterName := strings.Join(args[1:], " ")
	currentPlayedCharacter, exists := funk.Find(state.Players, func(p gp.Player) bool { return p.Id == maybePlayer.User.ID }).(gp.Player)
	if exists {
		currentPlayedCharacterChannel, err := s.GetChannel(sess, currentPlayedCharacter.Channel)
		if err != nil {
			return err
		}
		if currentPlayedCharacterChannel != nil {
			_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.player.have-already-character", map[string]interface{}{
				"Name":          maybePlayer.User.Username,
				"CharacterName": currentPlayedCharacter.CharacterName,
			})
			return err
		}
	}

	playerChannel, err := sess.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
		Name:     strings.ToLower(strings.Join(args[1:], "-")),
		ParentID: playersChannel.ID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{ID: state.Roles.Master, Type: discordgo.PermissionOverwriteTypeRole, Allow: discord.AllPermissionForChannel},
			{ID: maybePlayer.User.ID, Type: discordgo.PermissionOverwriteTypeMember, Allow: discord.AllPermissionForChannel},
			{ID: guild.ID, Deny: discordgo.PermissionViewChannel},
		},
	})

	err = sess.GuildMemberNickname(guild.ID, maybePlayer.User.ID, characterName)
	if err != nil {
		_, _ = discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.player.cannot-change-name", map[string]interface{}{
			"Name":          maybePlayer.User.Username,
			"CharacterName": currentPlayedCharacter.CharacterName,
		})
	}
	newPlayers := funk.Filter(state.Players, func(p gp.Player) bool { return p.Id != maybePlayer.User.ID }).([]gp.Player)
	if err := s.Save(state, Fields{"players": append(newPlayers, gp.Player{Id: maybePlayer.User.ID, CharacterName: characterName, Channel: playerChannel.ID})}, context.Background()); err != nil {
		return err
	}
	return sess.MessageReactionAdd(msg.ChannelID, msg.ID, "👌")
}
