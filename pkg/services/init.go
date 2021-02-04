package services

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/i18n"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
	"log"
)

func (s *Services) Init(msg *discordgo.Message, sess *discordgo.Session, state *models.ServerState) error {
	guild, err := sess.Guild(msg.GuildID)
	if err != nil {
		return err
	}
	log.Printf("init for guild '%s' (%s)", guild.Name, guild.ID)
	messages := make([]*discordgo.Message, 0)
	sendMessage := func(key string, props interface{}) error {
		mess, err := discord.SendMessage(msg.ChannelID, sess, state.Language, key, props)
		if err != nil {
			return err
		}
		messages = append(messages, mess)
		return nil
	}
	if msg.Author.ID != guild.OwnerID {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.init.unauthorized", nil)
		return err
	}
	if err := sendMessage("messages.init.start", state); err != nil {
		return err
	}

	masterRole := s.GetRole(guild, state.Roles.Players)
	if masterRole == nil {
		log.Printf(" - master role not exists, create it")
		masterRole, err = sess.GuildRoleCreate(guild.ID)
		if err != nil {
			return err
		}
		masterRole, err = sess.GuildRoleEdit(guild.ID, masterRole.ID, i18n.Must(state.Language, "name.role.master", nil), discord.GOLD, true, int(masterRole.Permissions), false)
		if err != nil {
			return err
		}
		state.Roles.Master = masterRole.ID
		if err := s.Save(state, Fields{"roles.master": masterRole.ID}, context.Background()); err != nil {
			return err
		}
		if err := sendMessage("messages.init.master-role-created", nil); err != nil {
			return err
		}
	}

	playersRole := s.GetRole(guild, state.Roles.Players)
	if playersRole == nil {
		log.Printf(" - players role not exists, create it")
		playersRole, err = sess.GuildRoleCreate(guild.ID)
		if err != nil {
			return err
		}
		playersRole, err = sess.GuildRoleEdit(guild.ID, playersRole.ID, i18n.Must(state.Language, "name.role.players", nil), discord.DARK_GREEN, true, int(playersRole.Permissions), false)
		if err != nil {
			return err
		}
		state.Roles.Players = playersRole.ID
		if err := s.Save(state, Fields{"roles.players": playersRole.ID}, context.Background()); err != nil {
			return err
		}
		if err := sendMessage("messages.init.players-role-created", nil); err != nil {
			return err
		}
	}

	playersChannel, err := s.GetChannel(sess, state.Channels.Players)
	if err != nil {
		return err
	}
	if playersChannel == nil {
		log.Printf("  - players channel not exists, create it")
		playersChannel, err = sess.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
			Name: i18n.Must(state.Language, "name.channel.players", nil),
			Type: discordgo.ChannelTypeGuildCategory,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{ID: state.Roles.Master, Type: discordgo.PermissionOverwriteTypeRole, Allow: discord.AllPermissionForChannel},
				{ID: guild.ID, Deny: discordgo.PermissionViewChannel},
			},
		})
		if err != nil {
			return err
		}
		if err := s.Save(state, Fields{"channels.players": playersChannel.ID}, context.Background()); err != nil {
			return err
		}
	}

	withBotChannel, err := s.GetChannel(sess, state.Channels.WithBot)
	if err != nil {
		return err
	}
	if withBotChannel == nil {
		log.Printf("  - system channel not exists, create it")
		withBotChannel, err = sess.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
			Name:     i18n.Must(state.Language, "name.channel.with-bot", nil),
			ParentID: playersChannel.ID,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{ID: state.Roles.Master, Type: discordgo.PermissionOverwriteTypeRole, Allow: discord.AllPermissionForChannel},
				{ID: guild.ID, Deny: discord.AllPermissionForChannelExceptView, Allow: discordgo.PermissionViewChannel},
			},
		})
		if err != nil {
			return err
		}
		if err := s.Save(state, Fields{"channels.withBot": withBotChannel.ID, "isInitialized": true}, context.Background()); err != nil {
			return err
		}
		_, err := discord.SendMessage(withBotChannel.ID, sess, state.Language, "messages.init.config-hint", nil)
		if err != nil {
			return err
		}
		_, err = discord.SendMessage(withBotChannel.ID, sess, state.Language, "messages.init.finished", nil)
		if err != nil {
			return err
		}
	}
	for _, m := range messages {
		_ = sess.ChannelMessageDelete(m.ChannelID, m.ID)
	}
	_ = sess.MessageReactionAdd(msg.ChannelID, msg.ID, "👌")

	return nil
}
