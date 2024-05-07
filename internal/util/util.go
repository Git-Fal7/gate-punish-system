package util

import (
	"strings"

	"github.com/git-fal7/gate-punish-system/internal/config"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// utils
func ReplaceAll(str string, replaceMap map[string]string) string {
	for key, value := range replaceMap {
		str = strings.ReplaceAll(str, key, value)
	}
	return str
}

func BroadcastPunishment(p *proxy.Proxy, component component.Component) {
	switch strings.ToUpper(config.ViperConfig.GetString("config.broadcast_punishment_to")) {
	case "ALL":
		go func() {
			for _, player := range p.Players() {
				player.SendMessage(component)
			}
		}()
	case "STAFF":
		go func() {
			for _, player := range p.Players() {
				if player.HasPermission(config.ViperConfig.GetString("permission.staff")) {
					player.SendMessage(component)
				}
			}
		}()
	case "NONE":
		break
	}
}
