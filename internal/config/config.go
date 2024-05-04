package config

import (
	"log"

	"github.com/spf13/viper"
)

var ViperConfig = viper.New()

func InitConfig() {
	ViperConfig.AddConfigPath(".")
	ViperConfig.SetConfigName("gatepunishsystem")
	ViperConfig.SetConfigType("yaml")

	ViperConfig.SetDefault("database.hostname", "localhost")
	ViperConfig.SetDefault("database.port", 5432)
	ViperConfig.SetDefault("database.username", "admin")
	ViperConfig.SetDefault("database.password", "adminpassword")
	ViperConfig.SetDefault("database.database", "punishsystem")

	ViperConfig.SetDefault("config.kick.permission", "git-fal7.punishsystem.kick")

	ViperConfig.SetDefault("messages.kick.format", "/kick (player) (reason)")
	ViperConfig.SetDefault("messages.kick.kick_message", "You got kicked\nReason: %reason%\nStaff: %staff%")

	ViperConfig.SetDefault("messages.error.playerNotFound", "Player not found")

	err := ViperConfig.ReadInConfig()
	if err != nil {
		// Create config file
		log.Println("Couldn't find gatepunishsystem.yml, creating a new config file")
		ViperConfig.WriteConfigAs("./gatepunishsystem.yml")
	}

}
