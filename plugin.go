package friendforgate

import (
	"github.com/git-fal7/gate-punish-system/internal/plugin"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Plugin = proxy.Plugin{
	Name: "Punish-system",
	Init: plugin.InitPlugin,
}
