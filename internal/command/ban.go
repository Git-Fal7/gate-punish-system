package command

import (
	"context"
	"strings"
	"time"

	"github.com/git-fal7/gate-punish-system/internal/config"
	"github.com/git-fal7/gate-punish-system/internal/database"
	"github.com/git-fal7/gate-punish-system/internal/timeutils"
	"github.com/git-fal7/gate-punish-system/internal/util"
	"github.com/google/uuid"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func banCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("ban").
		Requires(command.Requires(func(c *command.RequiresContext) bool {
			return c.Source.HasPermission(config.ViperConfig.GetString("permission.ban"))
		})).
		Executes(command.Command(func(c *command.Context) error {
			c.Source.SendMessage(&component.Text{
				Content: config.ViperConfig.GetString("messages.ban.format"),
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
					target := p.PlayerByName(targetName)
					if target == nil {
						c.Source.SendMessage(&component.Text{
							Content: config.ViperConfig.GetString("messages.error.playerNotFound"),
						})
						return nil
					}
					if target.HasPermission(config.ViperConfig.GetString("permission.staff")) {
						c.Source.SendMessage(&component.Text{
							Content: config.ViperConfig.GetString("messages.error.cantBanPlayer"),
						})
						return nil
					}
					return nil
				})).
				Then(
					brigodier.Argument("duration", brigodier.String).
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
									Content: config.ViperConfig.GetString("messages.error.cantBanPlayer"),
								})
								return nil
							}
							durationString := c.String("duration")
							duration, err := timeutils.ParseDuration(durationString)
							if err != nil || duration < time.Second {
								c.Source.SendMessage(&component.Text{
									Content: config.ViperConfig.GetString("messages.error.wrongDurationFormat"),
								})
								return nil
							}
							banPlayer(p, c.Source, target, "No Reason", duration)
							return nil
						})).
						Then(
							brigodier.Argument("reason", brigodier.String).
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
											Content: config.ViperConfig.GetString("messages.error.cantBanPlayer"),
										})
										return nil
									}
									durationString := c.String("duration")
									duration, err := timeutils.ParseDuration(durationString)
									if err != nil || duration < time.Second {
										c.Source.SendMessage(&component.Text{
											Content: config.ViperConfig.GetString("messages.error.wrongDurationFormat"),
										})
										return nil
									}
									reason := c.String("reason")
									banPlayer(p, c.Source, target, reason, duration)
									return nil
								},
								))),
				),
		)
}

func banPlayer(p *proxy.Proxy, source command.Source, target proxy.Player, reason string, duration time.Duration) {
	staffPlayer, ok := source.(proxy.Player)
	staffName := "Console"
	if ok {
		staffName = staffPlayer.Username()
	}
	timeEnds := time.Now().Add(duration)
	database.DB.PunishPlayer(context.Background(), database.PunishPlayerParams{
		UserUuid:   uuid.UUID(target.ID()),
		Reason:     reason,
		DoneBy:     staffName,
		PunishType: database.PunishtypeBAN,
		TimeEnds:   timeEnds,
	})
	target.Disconnect(&component.Text{
		Content: util.ReplaceAll(config.ViperConfig.GetString("messages.ban.ban_message"),
			map[string]string{
				"%target%": target.Username(),
				"%reason%": reason,
				"%staff%":  staffName,
				"%time%":   timeEnds.Format(config.ViperConfig.GetString("config.time_format")),
			}),
	})
	util.BroadcastPunishment(p, &component.Text{
		Content: util.ReplaceAll(config.ViperConfig.GetString("messages.ban.punish"),
			map[string]string{
				"%target%": target.Username(),
				"%reason%": reason,
				"%staff%":  staffName,
				"%time%":   timeEnds.Format(config.ViperConfig.GetString("config.time_format")),
			}),
	})
}
