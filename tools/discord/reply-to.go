package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/i18n"
	"log"
)

func ReplyTo(msg *discordgo.Message, sess *discordgo.Session, lang string, key string, args interface{}) (*discordgo.Message, error) {
	content, err := i18n.Translate(lang, key, args)
	if err != nil {
		log.Printf("failed to load i18n message %s", key)
		return nil, err
	}

	return ReplyToWithContent(msg, sess, content)
}

func ReplyToWithContent(msg *discordgo.Message, sess *discordgo.Session, content string) (*discordgo.Message, error) {
	return sess.ChannelMessageSendReply(msg.ChannelID, content, msg.Reference())
}
