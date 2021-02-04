package errors

import (
	"github.com/bwmarrin/discordgo"
)

func IsHttpError(err error, status int) bool {
	switch t := err.(type) {
	case discordgo.RESTError:
		return t.Response.StatusCode == status
	case *discordgo.RESTError:
		return t.Response.StatusCode == status
	}
	return false
}
