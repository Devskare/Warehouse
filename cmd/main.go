package main

import (
	"context"
	"fmt"
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
/*
	testRepo := repository.NewWHouseRepository(sqlDB)
	now := time.Now()
	test := models.ProductModel{
		Article:      1,
		ProductName:  "Test product",
		StorageID:    1,
		DeliveryDate: &now,
		ExpireDate:   &now,
		Weight:       10.5,
	}
	storageTest := models.StorageModel{
		MaxWeight: 100,
	}

	err = testRepo.StorageADD(ctx, storageTest.MaxWeight)
	if err != nil {
		log.Error("failed to add storage to database, fatal error ", err)
		panic(err)
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
	/*
		err = testRepo.ProductDelete(ctx, test.Article)
		if err != nil {
			log.Error("failed to delete product from database, fatal error ", err)
			panic(err)
		}

		test.Article = 2
		test.ProductName = "Test product 2"
		test.DeliveryDate = time.Now()
		test.ExpireDate = time.Now()
		test.Weight = 1000
		/*
			err = testRepo.ProductADD(ctx, test)
			if err != nil {
				log.Error("failed to add product to database, fatal error ", err)
			}


	testProducts, err := testRepo.ListProducts(ctx)
	if err != nil {
		log.Error("failed to list products from database, fatal error ", err)
		panic(err)
	}
	for _, product := range testProducts {
		fmt.Println(product)
	}

	TestStorages, err := testRepo.ListStorages(ctx)
	if err != nil {
		log.Error("failed to list storages from database, fatal error ", err)
		panic(err)
	}
	for _, storage := range TestStorages {
		fmt.Println(storage)
	}

	TestProduct, err := testRepo.GetProduct(ctx, testProducts[0].Article)
	if err != nil {
		log.Error("failed to get product from database, fatal error ", err)
		panic(err)
	}
	fmt.Println(TestProduct)
*/
 */

}
