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

	ViperConfig.SetDefault("permission.kick", "git-fal7.punishsystem.kick")
	ViperConfig.SetDefault("permission.mute", "git-fal7.punishsystem.mute")
	ViperConfig.SetDefault("permission.ban", "git-fal7.punishsystem.ban")
	ViperConfig.SetDefault("permission.unban", "git-fal7.punishsystem.unban")
	ViperConfig.SetDefault("permission.staff", "git-fal7.punishsystem.staff")

	// ALL - STAFF - NONE
	ViperConfig.SetDefault("config.broadcast_punishment_to", "ALL")
	ViperConfig.SetDefault("config.time_format", "2006-01-02 15:04:05")

	ViperConfig.SetDefault("messages.unban.format", "/unban (player)")
	ViperConfig.SetDefault("messages.unban.broadcast", "%target% has been unbanned by %staff%")

	ViperConfig.SetDefault("messages.ban.format", "/ban (player) (duration) (reason)")
	ViperConfig.SetDefault("messages.ban.broadcast", "%target% has been banned by %staff% for %reason% for %time%")
	ViperConfig.SetDefault("messages.ban.ban_message", "You got banned\nfor %reason%\nby %staff%\nEnds in %time%")

	ViperConfig.SetDefault("messages.mute.format", "/ban (player) (duration) (reason)")
	ViperConfig.SetDefault("messages.mute.broadcast", "%target% has been muted by %staff% for %reason% for %time%")
	ViperConfig.SetDefault("messages.mute.mute_message", "You got muted for %reason% by %staff% Ends in %time%")

	ViperConfig.SetDefault("messages.kick.format", "/kick (player) (reason)")
	ViperConfig.SetDefault("messages.kick.broadcast", "%target% has been banned by %staff% for %reason%")
	ViperConfig.SetDefault("messages.kick.kick_message", "You got kicked\nReason: %reason%\nStaff: %staff%")

	ViperConfig.SetDefault("messages.error.wrongDurationFormat", "Wrong duration time!, it should be as follows: 3h10m2s")
	ViperConfig.SetDefault("messages.error.cantBanPlayer", "Cannot ban this player")
	ViperConfig.SetDefault("messages.error.cantKickPlayer", "Cannot kick this player")
	ViperConfig.SetDefault("messages.error.cantMutePlayer", "Cannot mute this player")
	ViperConfig.SetDefault("messages.error.playerUnbanned", "Player is already unbanned")
	ViperConfig.SetDefault("messages.error.playerNotFound", "Player not found")

	err := ViperConfig.ReadInConfig()
	if err != nil {
		// Create config file
		log.Println("Couldn't find gatepunishsystem.yml, creating a new config file")
		ViperConfig.WriteConfigAs("./gatepunishsystem.yml")
	}

}
