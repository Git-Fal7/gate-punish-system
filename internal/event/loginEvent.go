package event

import (
	"context"
	"log"

	"github.com/git-fal7/gate-punish-system/internal/config"
	"github.com/git-fal7/gate-punish-system/internal/database"
	"github.com/git-fal7/gate-punish-system/internal/util"
	"github.com/google/uuid"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func loginEvent() func(*proxy.LoginEvent) {
	return func(e *proxy.LoginEvent) {
		player := e.Player()

		// log into lookup table
		logIntoLookupTableParam := database.LogIntoLookupTableParams{
			UserUuid: uuid.UUID(player.ID()),
			UserName: player.Username(),
		}
		err := database.DB.LogIntoLookupTable(context.Background(), logIntoLookupTableParam)
		if err != nil {
			log.Println(err)
		}

		// check if banned
		v, err := database.DB.IsPunishedPlayer(context.Background(), database.IsPunishedPlayerParams{
			UserUuid:   uuid.UUID(player.ID()),
			PunishType: database.PunishtypeBAN,
		})
		if err == nil {
			player.Disconnect(&component.Text{
				Content: util.ReplaceAll(config.ViperConfig.GetString("messages.ban.ban_message"),
					map[string]string{
						"%reason%": v.Reason,
						"%staff%":  v.DoneBy,
						"%time%":   v.TimeEnds.Format(config.ViperConfig.GetString("config.time_format")),
					}),
			})
		}
	}
}
