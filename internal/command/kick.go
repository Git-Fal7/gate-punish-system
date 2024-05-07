package command

import (
	"strings"

	"github.com/git-fal7/gate-punish-system/internal/config"
	"github.com/git-fal7/gate-punish-system/internal/util"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func kickCommand(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("kick").
		Requires(command.Requires(func(c *command.RequiresContext) bool {
			return c.Source.HasPermission(config.ViperConfig.GetString("permission.kick"))
		})).
		Executes(command.Command(func(c *command.Context) error {
			c.Source.SendMessage(&component.Text{
				Content: config.ViperConfig.GetString("messages.kick.format"),
			})
			return nil
		})).Then(
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
				target := p.PlayerByName(targetName)
				if target == nil {
					c.Source.SendMessage(&component.Text{
						Content: config.ViperConfig.GetString("messages.error.playerNotFound"),
					})
					return nil
				}
				if target.HasPermission(config.ViperConfig.GetString("permission.staff")) {
					c.Source.SendMessage(&component.Text{
						Content: config.ViperConfig.GetString("messages.error.cantKickPlayer"),
					})
					return nil
				}
				kickPlayer(p, target, "No reason", c.Source)
				return nil
			})).
			Then(brigodier.Argument("message", brigodier.StringPhrase).
				Executes(command.Command(func(c *command.Context) error {
					targetName := c.String("player")
					target := p.PlayerByName(targetName)
					if target == nil {
						c.Source.SendMessage(&component.Text{
							Content: config.ViperConfig.GetString("messages.error.playerNotFound"),
						})
						return nil
					}
					if target.HasPermission(config.ViperConfig.GetString("permission.staff")) {
						c.Source.SendMessage(&component.Text{
							Content: config.ViperConfig.GetString("messages.error.cantKickPlayer"),
						})
						return nil
					}
					message := c.String("message")
					kickPlayer(p, target, message, c.Source)
					return nil
				})),
			))
}

func kickPlayer(p *proxy.Proxy, target proxy.Player, reason string, source command.Source) {
	staffPlayer, ok := source.(proxy.Player)
	staffName := "Console"
	if ok {
		staffName = staffPlayer.Username()
	}
	target.Disconnect(&component.Text{
		Content: util.ReplaceAll(config.ViperConfig.GetString("messages.kick.kick_message"),
			map[string]string{
				"%target%": target.Username(),
				"%reason%": reason,
				"%staff%":  staffName,
			}),
	})
	util.BroadcastPunishment(p, &component.Text{
		Content: util.ReplaceAll(config.ViperConfig.GetString("messages.kick.punish"),
			map[string]string{
				"%target%": target.Username(),
				"%reason%": reason,
				"%staff%":  staffName,
			}),
	})
}
