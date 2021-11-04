package config

import (
	"github.com/spf13/viper"
)

func MongoUri() string     { return viper.GetString("MONGO_URI") }
func DiscordToken() string { return viper.GetString("DISCORD_TOKEN") }
func DiscordAppId() string { return viper.GetString("DISCORD_APPLICATION_ID") }
