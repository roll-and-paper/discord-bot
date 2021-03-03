package models

import "time"

type Player struct {
	Id            string `json:"id" bson:"id"`
	CharacterName string `json:"characterName" bson:"characterName"`
	Channel       string `json:"channel" bson:"channel"`
}

const (
	Fr = "fr"
	En = "en"
)

type Config struct {
	Prefix     string `json:"prefix" bson:"prefix"`
	GameSystem string `json:"gameSystem" bson:"gameSystem"`
}

type Roles struct {
	Master  string `json:"master,omitempty" bson:"master,omitempty"`
	Players string `json:"players,omitempty" bson:"players,omitempty"`
}

type Channels struct {
	Players string `json:"players,omitempty" bson:"players,omitempty"`
	WithBot string `json:"withBot,omitempty" bson:"withBot,omitempty"`
}

type ServerState struct {
	Id            string    `json:"id" bson:"_id"`
	Config        Config    `json:"config" bson:"config"`
	Roles         Roles     `json:"roles" bson:"roles"`
	Channels      Channels  `json:"channels" bson:"channels"`
	Master        string    `json:"master,omitempty" bson:"master,omitempty"`
	Players       []Player  `json:"players" bson:"players"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	IsInitialized bool      `json:"isInitialized" bson:"isInitialized"`
	Language      string    `json:"language" bson:"language"`
}

func NewServerState(id string) *ServerState {
	return &ServerState{
		Id: id,
		Config: Config{
			Prefix: "$",
		},
		Players:   make([]Player, 0),
		Language:  Fr,
		CreatedAt: time.Now(),
	}
}
