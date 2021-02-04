package discord

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

const AllPermissionForChannel int64 = discordgo.PermissionViewChannel |
	discordgo.PermissionAddReactions |
	discordgo.PermissionSendMessages |
	discordgo.PermissionSendTTSMessages |
	discordgo.PermissionManageMessages |
	discordgo.PermissionEmbedLinks |
	discordgo.PermissionAttachFiles |
	discordgo.PermissionReadMessageHistory |
	discordgo.PermissionUseExternalEmojis

func GetUserIdFromMention(mention string) string {
	if strings.HasPrefix(mention, "<@!") && strings.HasSuffix(mention, ">") {
		return strings.TrimSuffix(strings.TrimPrefix(mention, "<@!"), ">")
	}
	return ""
}
