package command

import (
	"github.com/git-fal7/gate-punish-system/internal/config"
	"github.com/git-fal7/gate-punish-system/internal/stringutil"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func kickCommand(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("kick").
		Requires(command.Requires(func(c *command.RequiresContext) bool {
			return c.Source.HasPermission(config.ViperConfig.GetString("config.kick.permission"))
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
				for _, target := range p.Players() {
					b.Suggest(target.Username())
				}
				return b.Build()
			})).
			Executes(command.Command(func(c *command.Context) error {
				playerStr := c.String("player")
				plr := p.PlayerByName(playerStr)
				if plr == nil {
					c.Source.SendMessage(&component.Text{
						Content: config.ViperConfig.GetString("messages.error.playerNotFound"),
					})
					return nil
				}
				kickPlayer(plr, "No reason", c.Source)
				return nil
			})).
			Then(brigodier.Argument("message", brigodier.StringPhrase).
				Executes(command.Command(func(c *command.Context) error {
					playerStr := c.String("player")
					plr := p.PlayerByName(playerStr)
					if plr == nil {
						c.Source.SendMessage(&component.Text{
							Content: config.ViperConfig.GetString("messages.error.playerNotFound"),
						})
						return nil
					}
					message := c.String("message")
					kickPlayer(plr, message, c.Source)
					return nil
				})),
			))
}

func kickPlayer(player proxy.Player, reason string, source command.Source) {
	staffPlayer, ok := source.(proxy.Player)
	staffName := "Console"
	if ok {
		staffName = staffPlayer.Username()
	}
	player.Disconnect(&component.Text{
		Content: stringutil.ReplaceAll(config.ViperConfig.GetString("messages.kick.kick_message"),
			map[string]string{
				"%reason%": reason,
				"%staff%":  staffName,
			}),
	})
}
