package services

import "github.com/dohr-michael/roll-and-paper-bot/tools/storage"

type Services struct{}

func NewServices(col storage.Storage) *Services {
	return &Services{}
}
