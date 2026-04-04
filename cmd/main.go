package main

import (
	"context"
	"time"
	"warehouse/config"
	"warehouse/logger"
	"warehouse/modules/Whouse/models"
	"warehouse/modules/Whouse/repository"
	"warehouse/modules/db"
)

func main() {

	time.Sleep(5 * time.Second)

	appConf := config.MustLoadConfig()

	log := logger.Initlogger(appConf.LogLevel, appConf.Production)
	log.Info("info")

	sqlDB, err := db.NewSqlDB(log, &appConf.DB)
	if err != nil {
		log.Error("failed to connect to database, fatal error ", err)
		panic(err)
	}
	ctx := context.Background()

	//if err := createSQL.CreateTables(ctx, conn); err != nil { panic(err) }

	testRepo := repository.NewWHouseRepository(sqlDB)
	test := models.ProductModel{
		Article:      1,
		ProductName:  "Test product",
		StorageID:    2,
		DeliveryDate: time.Now(),
		ExpireDate:   time.Now(),
		Weight:       10.5,
	}

	err = testRepo.ProductADD(ctx, test)
	if err != nil {
		log.Error("failed to add product to database, fatal error ", err)
		panic(err)
	}
	err = testRepo.ProductUpdate(ctx, test)
	if err != nil {
		log.Error("failed to update product to database, fatal error ", err)
		panic(err)
	}
	err = testRepo.ProductDelete(ctx, test.Article)
	if err != nil {
		log.Error("failed to delete product from repository, fatal error ", err)
		panic(err)
	}

}
