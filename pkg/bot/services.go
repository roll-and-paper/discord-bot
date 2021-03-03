package bot

import (
	"context"
	"github.com/bwmarrin/discordgo"
	gp "github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/tools/errors"
	"github.com/dohr-michael/roll-and-paper-bot/tools/storage"
	"github.com/thoas/go-funk"
)

type Fields map[string]interface{}

type Services struct {
	Save func(*gp.ServerState, Fields, context.Context) error
}

func (s *Services) GetRole(guild *discordgo.Guild, id string) *discordgo.Role {
	res, ok := funk.Find(guild.Roles, func(role *discordgo.Role) bool { return role.ID == id }).(*discordgo.Role)
	if !ok {
		return nil
	}
	return res
}

func (s *Services) GetChannel(sess *discordgo.Session, id string) (*discordgo.Channel, error) {
	if id != "" {
		channel, err := sess.Channel(id)
		if err != nil && errors.IsHttpError(err, 404) {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		return channel, nil
	}
	return nil, nil
}

func New(col storage.Storage) *Services {
	return &Services{
		Save: func(state *gp.ServerState, toUpdate Fields, ctx context.Context) error {
			return col.Save(state.Id, toUpdate, state, ctx)
		},
	}
}
