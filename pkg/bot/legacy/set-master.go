package legacy

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
	"github.com/thoas/go-funk"
)

func setMaster(s *Services, msg *discordgo.Message, sess *discordgo.Session, state *models.ServerState, args []string) error {
	if len(args) == 0 {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.master.missing", state)
		return err
	}
	guild, err := sess.Guild(msg.GuildID)
	if err != nil {
		return err
	}

	masterRole := s.GetRole(guild, state.Roles.Master)
	if masterRole == nil {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.master.role-missing", state)
		return err
	}
	maybeMasterId := discord.GetUserIdFromMention(args[0])
	if maybeMasterId == "" {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.not-a-person", map[string]interface{}{"Name": args[0]})
		return err
	}
	var currentMaster *discordgo.Member
	if state.Master != "" {
		currentMaster, _ = sess.GuildMember(msg.GuildID, state.Master)
	}
	if currentMaster != nil && currentMaster.User.ID == maybeMasterId && funk.Contains(currentMaster.Roles, masterRole.ID) {
		return sess.MessageReactionAdd(msg.ChannelID, msg.ID, "👌")
	}

	maybeMaster, err := sess.GuildMember(msg.GuildID, maybeMasterId)
	if err != nil || maybeMaster == nil {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.unknown-person", map[string]interface{}{"Name": maybeMasterId})
		return err
	}
	if currentMaster != nil {
		_ = sess.GuildMemberRoleRemove(msg.GuildID, currentMaster.User.ID, masterRole.ID)
	}
	if err := sess.GuildMemberRoleAdd(msg.GuildID, maybeMaster.User.ID, masterRole.ID); err != nil {
		return err
	}
	if err := s.Save(state, Fields{"master": maybeMaster.User.ID}, context.Background()); err != nil {
		return err
	}
	return sess.MessageReactionAdd(msg.ChannelID, msg.ID, "👌")
}
