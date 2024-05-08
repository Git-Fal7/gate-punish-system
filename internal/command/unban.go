package command

import (
	"context"
	"strings"

	"github.com/git-fal7/gate-punish-system/internal/config"
	"github.com/git-fal7/gate-punish-system/internal/database"
	"github.com/git-fal7/gate-punish-system/internal/util"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func unbanCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("unban").
		Requires(command.Requires(func(c *command.RequiresContext) bool {
			return c.Source.HasPermission(config.ViperConfig.GetString("permission.unban"))
		})).
		Executes(command.Command(func(c *command.Context) error {
			c.Source.SendMessage(&component.Text{
				Content: config.ViperConfig.GetString("messages.unban.format"),
			})
			return nil
		})).
		Then(
			brigodier.Argument("player", brigodier.String).
				Suggests(command.SuggestFunc(func(
					c *command.Context,
					b *brigodier.SuggestionsBuilder,
				) *brigodier.Suggestions {
					arg := strings.ToLower(c.String("player"))
					for _, target := range p.Players() {
						if strings.HasPrefix(strings.ToLower(target.Username()), arg) {
							if !target.HasPermission(config.ViperConfig.GetString("permission.staff")) {
								b.Suggest(target.Username())
							}
						}
					}
					return b.Build()
				})).
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
						PunishType: database.PunishtypeBAN,
					})
					if err != nil {
						c.Source.SendMessage(&component.Text{
							Content: config.ViperConfig.GetString("messages.error.playerUnbanned"),
						})
						return nil
					}
					// broadcast that player is not
					staffPlayer, ok := c.Source.(proxy.Player)
					staffName := "Console"
					if ok {
						staffName = staffPlayer.Username()
					}
					util.BroadcastPunishment(p, &component.Text{
						Content: util.ReplaceAll(
							config.ViperConfig.GetString("messages.unban.broadcast"),
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
