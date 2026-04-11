package main

import (
	"time"
	"warehouse"
	"warehouse/config"
	"warehouse/logger"
)

func main() {

	time.Sleep(2 * time.Second)

	appConf := config.MustLoadConfig(".env.test")

	log := logger.Initlogger(appConf.LogLevel, appConf.Production)
	log.Info("info")

	warehouse.TestMain(log)

}
