package command

import (
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func Init(p *proxy.Proxy) {
	p.Command().Register(kickCommand(p))
	p.Command().Register(banCmd(p))
	p.Command().Register(muteCmd(p))
	p.Command().Register(unbanCmd(p))
}
