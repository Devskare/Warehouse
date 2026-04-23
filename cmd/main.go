package main

import (
	"log/slog"
	"time"
	"warehouse/config"
	"warehouse/logger"
	"warehouse/modules/db"
)

func main() {

	time.Sleep(2 * time.Second)

	appConf := config.MustLoadConfig(".env")
	log := logger.Initlogger(appConf.LogLevel, appConf.Production)

	_, err := db.NewSqlDB(log, &appConf.DB)
	if err != nil {
		log.Error("failed to connect sql data base", slog.Any("error", err))
	}
}
