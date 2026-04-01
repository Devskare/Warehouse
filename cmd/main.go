package main

import (
	"context"
	"time"
	"warehouse/config"
	"warehouse/logger"
	"warehouse/modules/db"
	"warehouse/modules/db/createSQL"
)

func main() {

	time.Sleep(5 * time.Second)

	appConf := config.MustLoadConfig()

	ctx := context.Background()
	conn, err := db.Connection(ctx, appConf)
	if err != nil {
		panic(err)
	}
	if err := createSQL.CreateTables(ctx, conn); err != nil {
		panic(err)

	}
	log := logger.Initlogger(appConf.LogLevel, appConf.Production)
	log.Info("info")
}
