package plugin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/git-fal7/gate-punish-system/internal/command"
	"github.com/git-fal7/gate-punish-system/internal/config"
	"github.com/git-fal7/gate-punish-system/internal/database"
	"github.com/git-fal7/gate-punish-system/internal/event"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func InitPlugin(ctx context.Context, proxy *proxy.Proxy) error {
	config.InitConfig()
	db, err := sql.Open("pgx", fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s",
		config.ViperConfig.GetString("database.username"), config.ViperConfig.GetString("database.password"),
		config.ViperConfig.GetString("database.hostname"), config.ViperConfig.GetString("database.port"),
		config.ViperConfig.GetString("database.database")))
	if err != nil {
		return err
	}
	database.DB = database.New(db)
	// Init tables
	err = database.DB.CreatePunishType(context.Background())
	if err != nil {
		return err
	}
	err = database.DB.CreatePunishUsersTable(context.Background())
	if err != nil {
		return err
	}
	err = database.DB.CreateLookupUserTable(context.Background())
	if err != nil {
		return err
	}
	command.Init(proxy)
	event.Init(proxy)

	return nil
}
