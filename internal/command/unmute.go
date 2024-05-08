package command

import (
	"context"

	"github.com/git-fal7/gate-punish-system/internal/config"
	"github.com/git-fal7/gate-punish-system/internal/database"
	"github.com/git-fal7/gate-punish-system/internal/util"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func unmuteCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("unmute").
		Requires(command.Requires(func(c *command.RequiresContext) bool {
			return c.Source.HasPermission(config.ViperConfig.GetString("permission.unmute"))
		})).
		Executes(command.Command(func(c *command.Context) error {
			c.Source.SendMessage(&component.Text{
				Content: config.ViperConfig.GetString("messages.unmute.format"),
			})
			return nil
		})).
		Then(
			brigodier.Argument("player", brigodier.String).
				Executes(command.Command(func(c *command.Context) error {
					targetName := c.String("player")
					targetUUID, err := database.DB.GetPlayerUUID(context.Background(), targetName)
					if err != nil {
						c.Source.SendMessage(&component.Text{
							Content: config.ViperConfig.GetString("messages.error.playerNotFound"),
						})
						return nil
					}
					err = database.DB.UnpunishPlayer(context.Background(), database.UnpunishPlayerParams{
						UserUuid:   targetUUID,
						PunishType: database.PunishtypeMUTE,
					})
					if err != nil {
						c.Source.SendMessage(&component.Text{
							Content: config.ViperConfig.GetString("messages.error.playerUnmuted"),
						})
						return nil
					}
					staffPlayer, ok := c.Source.(proxy.Player)
					staffName := "Console"
					if ok {
						staffName = staffPlayer.Username()
					}
					util.BroadcastPunishment(p, &component.Text{
						Content: util.ReplaceAll(
							config.ViperConfig.GetString("messages.unmute.broadcast"),
							map[string]string{
								"%target%": targetName,
								"%staff%":  staffName,
							},
						),
					})
					return nil
				})),
		)
}
