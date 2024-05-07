package event

import (
	"context"

	"github.com/git-fal7/gate-punish-system/internal/config"
	"github.com/git-fal7/gate-punish-system/internal/database"
	"github.com/git-fal7/gate-punish-system/internal/util"
	"github.com/google/uuid"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func playerChatEvent() func(*proxy.PlayerChatEvent) {
	return func(e *proxy.PlayerChatEvent) {
		player := e.Player()
		v, err := database.DB.IsPunishedPlayer(context.Background(), database.IsPunishedPlayerParams{
			UserUuid:   uuid.UUID(player.ID()),
			PunishType: database.PunishtypeMUTE,
		})
		if err == nil {
			e.SetAllowed(false)
			player.SendMessage(&component.Text{
				Content: util.ReplaceAll(config.ViperConfig.GetString("messages.mute.mute_message"),
					map[string]string{
						"%reason%": v.Reason,
						"%staff%":  v.DoneBy,
						"%time%":   v.TimeEnds.Format(config.ViperConfig.GetString("config.time_format")),
					})},
			)
		}
	}
}
