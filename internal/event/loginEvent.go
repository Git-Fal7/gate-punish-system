package event

import (
	"context"
	"log"

	"github.com/git-fal7/gate-punish-system/internal/database"
	"github.com/google/uuid"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func loginEvent(p *proxy.Proxy) func(*proxy.LoginEvent) {
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

	}
}
