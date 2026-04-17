//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"
	"time"
	"warehouse/modules/Whouse/models"

	"warehouse/integration"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestAllRepo(t *testing.T) {
	allRepoTests(t)
}

func allRepoTests(t *testing.T) {
	testDB := integration.StartTestDB(t)

	db := testDB.DB
	ctx := testDB.Ctx
	t.Run("StorageAdd", func(t *testing.T) {
		testStorageADD(t, db, ctx)
	})

	t.Run("ProductAdd", func(t *testing.T) {
		testProductADD(t, db, ctx)
	})

	t.Run("ListStorages", func(t *testing.T) {
		testListStorages(t, db, ctx)
	})

	t.Run("ListProducts", func(t *testing.T) {
		testListProducts(t, db, ctx)
	})

	t.Run("GetProduct", func(t *testing.T) {
		testGetProduct(t, db, ctx)
	})

	t.Run("ProductUpdate", func(t *testing.T) {
		testProductUpdate(t, db, ctx)
	})

	t.Run("ProductDelete", func(t *testing.T) {
		testProductDelete(t, db, ctx)
	})

	t.Run("ProductExpire", func(t *testing.T) {
		testProductExpire(t, db, ctx)
	})

}

func cleanDB(t *testing.T, db *sqlx.DB) {
	t.Helper()

	_, err := db.Exec(`
        TRUNCATE TABLE products, storages RESTART IDENTITY CASCADE;
    `)
	require.NoError(t, err)
}

func testStorageADD(t *testing.T, db *sqlx.DB, ctx context.Context) {
	cleanDB(t, db)
	defer cleanDB(t, db)
	var storage models.StorageModel
	storage.MaxWeight = 100
	storage.CurrentWeight = 0

	repo := NewWHouseRepository(db)

	err := repo.StorageADD(ctx, storage.MaxWeight)
	require.NoError(t, err)

	var storageCheck models.StorageModel
	err = db.GetContext(ctx, &storageCheck, `
		SELECT id, max_weight, current_weight
		FROM storages
		WHERE id = 1
	`)
	require.NoError(t, err)

	require.InDelta(t, 100.0, storageCheck.MaxWeight, 0.001)
	require.InDelta(t, 0.0, storageCheck.CurrentWeight, 0.001)
}

func testProductADD(t *testing.T, db *sqlx.DB, ctx context.Context) {
	//подготавливаем дб и продукт.
	///////////////////////////////////////
	cleanDB(t, db)
	defer cleanDB(t, db)
	now := time.Now().UTC()
	storageID := 1
	product := models.ProductModel{
		Article:      1,
		ProductName:  "prod1",
		StorageID:    &storageID,
		DeliveryDate: &now,
		ExpireDate:   nil,
		Weight:       1,
	}
	repo := NewWHouseRepository(db)

	err := repo.StorageADD(ctx, 100)
	require.NoError(t, err)
	// тестим функцию, сверяем результаты.
	///////////////////////////////////////////////////////////////////////////////////////////////////
	err = repo.ProductADD(ctx, product)
	require.NoError(t, err)

	var dbProduct models.ProductModel

	err = db.GetContext(ctx, &dbProduct, `
		SELECT id, article, product_name, storage_id, delivery_date, expire_date, weight
		FROM products
		WHERE article = $1
	`, product.Article)

	require.NoError(t, err)

	require.Equal(t, product.Article, dbProduct.Article)
	require.Equal(t, product.ProductName, dbProduct.ProductName)
	require.NotNil(t, dbProduct.StorageID)
	require.Equal(t, *product.StorageID, *dbProduct.StorageID)
	require.InDelta(t, product.Weight, dbProduct.Weight, 0.001)

	require.WithinDuration(t,
		*product.DeliveryDate,
		*dbProduct.DeliveryDate,
		time.Second,
	)

	require.Nil(t, dbProduct.ExpireDate)

	var storage models.StorageModel
	err = db.GetContext(ctx, &storage, `
		SELECT id, current_weight
		FROM storages
		WHERE id = 1
	`)
	require.NoError(t, err)

	require.InDelta(t, 1.0, storage.CurrentWeight, 0.001)

}

func testListProducts(t *testing.T, db *sqlx.DB, ctx context.Context) {
	//Подготовка ДБ и Структур для теста.
	////////////////////////////////////
	cleanDB(t, db)
	defer cleanDB(t, db)

	now := time.Now().UTC()

	repo := NewWHouseRepository(db)
	storageID := 1
	product1 := models.ProductModel{
		Article:      1,
		ProductName:  "prod1",
		StorageID:    &storageID,
		DeliveryDate: &now,
		ExpireDate:   nil,
		Weight:       1.0,
	}

	product2 := models.ProductModel{
		Article:      2,
		ProductName:  "prod2",
		StorageID:    &storageID,
		DeliveryDate: &now,
		ExpireDate:   nil,
		Weight:       2.0,
	}
	require.NoError(t, repo.StorageADD(ctx, 100))
	require.NoError(t, repo.ProductADD(ctx, product1))
	require.NoError(t, repo.ProductADD(ctx, product2))

	//тест нужной функции.
	/////////////////////////////////////////////
	products, err := repo.ListProducts(ctx)
	require.NoError(t, err)

	require.Len(t, products, 2)

	//делаем мапу и сверяем результат.
	////////////////////////////////////////////////////
	got := make(map[int]models.ProductModel)
	for _, p := range products {
		got[p.Article] = p
	}

	p1, ok := got[product1.Article]
	require.True(t, ok)

	require.Equal(t, product1.ProductName, p1.ProductName)
	require.NotNil(t, p1.StorageID)
	require.Equal(t, *product1.StorageID, *p1.StorageID)
	require.InDelta(t, product1.Weight, p1.Weight, 0.001)
	require.NotNil(t, p1.DeliveryDate)
	require.Nil(t, p1.ExpireDate)

	p2, ok := got[product2.Article]
	require.True(t, ok)

	require.Equal(t, product2.ProductName, p2.ProductName)
	require.NotNil(t, p1.StorageID)
	require.Equal(t, *product2.StorageID, *p2.StorageID)
	require.InDelta(t, product2.Weight, p2.Weight, 0.001)
	require.NotNil(t, p2.DeliveryDate)
	require.Nil(t, p2.ExpireDate)
}

func testListStorages(t *testing.T, db *sqlx.DB, ctx context.Context) {
	//подготавливаем БД и добавляем туда storage1 и storage 2
	////////////////////////////////////////////////////////////////////////
	cleanDB(t, db)
	defer cleanDB(t, db)

	repo := NewWHouseRepository(db)

	// 🔹 ARRANGE
	storage1 := models.StorageModel{
		MaxWeight:     100.0,
		CurrentWeight: 0.0,
	}

	storage2 := models.StorageModel{
		MaxWeight:     200.0,
		CurrentWeight: 0.0,
	}

	err := repo.StorageADD(ctx, storage1.MaxWeight)
	require.NoError(t, err)

	err = repo.StorageADD(ctx, storage2.MaxWeight)
	require.NoError(t, err)
	//тестируем функцию и создаем мапу для сверки данных.
	/////////////////////////////////////////////////////////////////////////////////
	storages, err := repo.ListStorages(ctx)
	require.NoError(t, err)

	require.Len(t, storages, 2)

	got := make(map[int]models.StorageModel)
	for _, s := range storages {
		got[s.ID] = s
	}

	//сверяем ответ.
	/////////////////////////////////////////////////////////////////////////////////
	s1, ok := got[1]
	require.True(t, ok)

	require.InDelta(t, storage1.MaxWeight, s1.MaxWeight, 0.001)
	require.InDelta(t, storage1.CurrentWeight, s1.CurrentWeight, 0.001)

	s2, ok := got[2]
	require.True(t, ok)

	require.InDelta(t, storage2.MaxWeight, s2.MaxWeight, 0.001)
	require.InDelta(t, storage2.CurrentWeight, s2.CurrentWeight, 0.001)
}

func testGetProduct(t *testing.T, db *sqlx.DB, ctx context.Context) {
	//подготавливаем базу и продукт
	//////////////////////////////////////////////////////////////////////////
	cleanDB(t, db)
	defer cleanDB(t, db)

	repo := NewWHouseRepository(db)

	nuw := time.Now()
	storageID := 1
	product := models.ProductModel{
		Article:      1,
		ProductName:  "prod1",
		StorageID:    &storageID,
		DeliveryDate: &nuw,
		ExpireDate:   nil,
		Weight:       1.0,
	}
	err := repo.StorageADD(ctx, 100)
	require.NoError(t, err)
	err = repo.ProductADD(ctx, product)
	require.NoError(t, err)
	//тестируем функцию и сверяем значения
	///////////////////////////////////////////////////////////////////////////////////////////
	got, err := repo.GetProduct(ctx, product.Article)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.Equal(t, product.Article, got.Article)
	require.Equal(t, product.ProductName, got.ProductName)
	require.NotNil(t, got.StorageID)
	require.Equal(t, *product.StorageID, *got.StorageID)
	require.InDelta(t, product.Weight, got.Weight, 0.001)

	require.NotNil(t, got.DeliveryDate)
	require.WithinDuration(t,
		*product.DeliveryDate,
		*got.DeliveryDate,
		time.Second,
	)

	require.Nil(t, got.ExpireDate)
	//чистим дб и проверяем кейс, когда в дб нет продукта.
	//////////////////////////////////////////////////////////////////////////////
	cleanDB(t, db)

	got, err = repo.GetProduct(ctx, product.Article)
	require.Error(t, err)
	require.Nil(t, got)
}

func testProductUpdate(t *testing.T, db *sqlx.DB, ctx context.Context) {
	//Подготавливаем базу и продукты.
	/////////////////////////////////////////////////////////////////////////////////////////
	cleanDB(t, db)
	defer cleanDB(t, db)

	repo := NewWHouseRepository(db)

	now := time.Now().UTC()
	storageID1 := 1
	storageID2 := 2
	product := models.ProductModel{
		Article:      1,
		ProductName:  "prod1",
		StorageID:    &storageID1,
		DeliveryDate: &now,
		ExpireDate:   nil,
		Weight:       10,
	}

	err := repo.StorageADD(ctx, 100)
	require.NoError(t, err)
	err = repo.StorageADD(ctx, 100)
	require.NoError(t, err)
	err = repo.ProductADD(ctx, product)
	require.NoError(t, err)

	updatedProduct := product
	updatedProduct.StorageID = &storageID2
	updatedProduct.Weight = 15.0

	//Запускаем апдейт и проверяем новые значения.
	////////////////////////////////////////////////////////////////////////////////////////////////////////////
	err = repo.ProductUpdate(ctx, updatedProduct)
	require.NoError(t, err)
	got, err := repo.GetProduct(ctx, updatedProduct.Article)
	require.NoError(t, err)

	require.Equal(t, *updatedProduct.StorageID, *got.StorageID)
	require.InDelta(t, updatedProduct.Weight, got.Weight, 0.001)
	require.Equal(t, updatedProduct.ProductName, got.ProductName)

	var storage1 models.StorageModel
	err = db.GetContext(ctx, &storage1, `SELECT id, current_weight FROM storages WHERE id=$1`, product.StorageID)
	require.NoError(t, err)

	require.InDelta(t, 0.0, storage1.CurrentWeight, 0.001)

	var storage2 models.StorageModel
	err = db.GetContext(ctx, &storage2, `SELECT id, current_weight FROM storages WHERE id=$1`, updatedProduct.StorageID)
	require.NoError(t, err)

	require.InDelta(t, 15, storage2.CurrentWeight, 0.001)

}

func testProductDelete(t *testing.T, db *sqlx.DB, ctx context.Context) {
	//Подготавливаем базу.
	////////////////////////////////////////////////////////////////////////
	cleanDB(t, db)
	defer cleanDB(t, db)

	repo := NewWHouseRepository(db)

	now := time.Now().UTC()
	storageID1 := 1
	product := models.ProductModel{
		Article:      1,
		ProductName:  "prod1",
		StorageID:    &storageID1,
		DeliveryDate: &now,
		ExpireDate:   nil,
		Weight:       10,
	}

	err := repo.StorageADD(ctx, 100)
	require.NoError(t, err)
	err = repo.ProductADD(ctx, product)
	require.NoError(t, err)
	//тестируем функцию и новые значения
	//////////////////////////////////////////////////////////////////////////////////
	err = repo.ProductDelete(ctx, product.Article)
	require.NoError(t, err)

	got, err := repo.GetProduct(ctx, product.Article)
	require.Error(t, err)
	require.Nil(t, got)

	var storage models.StorageModel
	err = db.GetContext(ctx, &storage, `
		SELECT id, current_weight 
		FROM storages 
		WHERE id = $1
	`, product.StorageID)
	require.NoError(t, err)

	require.InDelta(t, 0, storage.CurrentWeight, 0.001)
}

func testProductExpire(t *testing.T, db *sqlx.DB, ctx context.Context) {
	cleanDB(t, db)
	defer cleanDB(t, db)

	repo := NewWHouseRepository(db)

	now := time.Now().UTC()

	err := repo.StorageADD(ctx, 100)
	require.NoError(t, err)
	storageID := 1
	product := models.ProductModel{
		Article:      1,
		ProductName:  "prod1",
		StorageID:    &storageID,
		DeliveryDate: &now,
		ExpireDate:   &now,
		Weight:       10,
	}

	err = repo.ProductADD(ctx, product)
	require.NoError(t, err)

	err = repo.ProductExpire(ctx, product)
	require.NoError(t, err)

	got, err := repo.GetProduct(ctx, product.Article)
	require.NoError(t, err)

	require.Nil(t, got.StorageID)

	require.NotNil(t, got.ExpireDate)

	require.WithinDuration(
		t,
		*product.ExpireDate,
		*got.ExpireDate,
		time.Second,
	)
}
