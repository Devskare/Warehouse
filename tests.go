package warehouse

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"warehouse/config"
	"warehouse/modules/Whouse/models"
	"warehouse/modules/Whouse/repository"
	"warehouse/modules/db"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var testDB *sqlx.DB

func TestMain(log *slog.Logger) {
	cfg := config.MustLoadConfig(".env.test")
	ctx := context.Background()

	dbConn, err := db.NewSqlDB(log, &cfg.DB)
	if err != nil {
		log.Error("failed to connect db: %v", err)
	}

	testDB = dbConn

	var m *migrate.Migrate
	m, err = runMigrations(&cfg.DB)
	if err != nil {
		log.Error("failed to run(init) migrations: %v", err)
	}
	err = m.Up()
	if err != nil {
		log.Error("failed to UP migrations: %v", err)
	}
	log.Info("migrations done")

	testRepo := repository.NewWHouseRepository(testDB)
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

	}
	err = testRepo.StorageADD(ctx, storageTest.MaxWeight)
	if err != nil {
		log.Error("failed to add storage to database, fatal error ", err)

	}
	log.Info("added storage to database")
	err = testRepo.ProductADD(ctx, test)
	if err != nil {
		log.Error("failed to add product to database, fatal error ", err)

	}
	log.Info("added product to database")

	test.ProductName = "Product Test"
	test.StorageID = 2

	err = testRepo.ProductUpdate(ctx, test)
	if err != nil {
		log.Error("failed to update product to database, fatal error ", err)
	}
	testProducts, err := testRepo.ListProducts(ctx)
	if err != nil {
		log.Error("failed to list products from database, fatal error ", err)

	}
	for _, product := range testProducts {
		fmt.Println(product)
	}

	TestStorages, err := testRepo.ListStorages(ctx)
	if err != nil {
		log.Error("failed to list storages from database, fatal error ", err)

	}
	for _, storage := range TestStorages {
		fmt.Println(storage)
	}

	TestProduct, err := testRepo.GetProduct(ctx, testProducts[0].Article)
	if err != nil {
		log.Error("failed to get product from database, fatal error ", err)

	}
	fmt.Println(TestProduct)

	err = testRepo.ProductDelete(ctx, test.Article)
	if err != nil {
		log.Error("failed to delete product from database, fatal error ", err)

	}

	err = m.Drop()
	if err != nil {
		log.Error("failed to DROP migrations: %v", err)
	}
	log.Info("dropped migrations")

	err = clearTables(testDB)
	if err != nil {
		log.Error("failed to clear tables: %v", err)
	}

}

func runMigrations(cfg *config.DB) (*migrate.Migrate, error) {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return nil, err
	}
	return m, err
}

func clearTables(db *sqlx.DB) error {
	_, err := db.Exec(`
        TRUNCATE TABLE products, storages RESTART IDENTITY CASCADE;
    `)
	return err
}
