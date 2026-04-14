//go:build integration
// +build integration

package repository

import (
	"fmt"
	"testing"
	"time"
	"warehouse/modules/Whouse/models"

	"warehouse/integration"

	"github.com/stretchr/testify/require"
)

func TestRepo(t *testing.T) {
	testDB := integration.StartTestDB(t)

	testRepo := NewWHouseRepository(testDB.DB)
	now := time.Now()

	testProduct := models.ProductModel{
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

	err := testRepo.StorageADD(testDB.Ctx, storageTest.MaxWeight)
	require.NoError(t, err)

	err = testRepo.StorageADD(testDB.Ctx, storageTest.MaxWeight)
	require.NoError(t, err)

	err = testRepo.ProductADD(testDB.Ctx, testProduct)
	require.NoError(t, err)

	testListProducts, err := testRepo.ListProducts(testDB.Ctx)
	fmt.Println(testListProducts)
	require.NoError(t, err)

	testListStorages, err := testRepo.ListStorages(testDB.Ctx)
	fmt.Println(testListStorages)
	require.NoError(t, err)

	TestGetProduct, err := testRepo.GetProduct(testDB.Ctx, testProduct.Article)
	fmt.Println(TestGetProduct)
	require.NoError(t, err)

	testProduct.StorageID = 2

	err = testRepo.ProductUpdate(testDB.Ctx, testProduct)
	require.NoError(t, err)

	err = testRepo.ProductDelete(testDB.Ctx, testProduct.Article)
	require.NoError(t, err)
}
