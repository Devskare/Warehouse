package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/devskar/warehouse/mocks"
	"github.com/devskar/warehouse/modules/Whouse/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProductADD_NotValid(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)

	service := NewProductService(repo, slog.Default())

	p := models.ProductModel{}

	err := service.ProductADD(context.Background(), p)

	require.Error(t, err)

	repo.AssertNotCalled(t, "ProductADD", mock.Anything, mock.Anything)
}

func TestProductADD_StorageIDNil(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())
	p := models.ProductModel{
		Article:     1,
		ProductName: "test",
		StorageID:   nil,
		Weight:      10,
	}
	err := service.ProductADD(context.Background(), p)
	require.Error(t, err)

	repo.AssertNotCalled(t, "ProductADD", mock.Anything, mock.Anything)
}

func TestProductADD_RepoError(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())
	storageID := 1
	p := models.ProductModel{
		Article:     1,
		ProductName: "test",
		StorageID:   &storageID,
		Weight:      10,
	}

	repo.On("ProductADD", mock.Anything, mock.AnythingOfType("models.ProductModel")).Return(errors.New("db error"))

	err := service.ProductADD(context.Background(), p)
	require.Error(t, err)
	require.EqualError(t, err, "db error")

	repo.AssertCalled(t, "ProductADD", mock.Anything, mock.AnythingOfType("models.ProductModel"))
}

func TestProductADD_Success(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	storageID := 1

	p := models.ProductModel{
		Article:     1,
		ProductName: "test",
		Weight:      10,
		StorageID:   &storageID,
	}

	repo.On("ProductADD", mock.Anything, mock.MatchedBy(func(p models.ProductModel) bool {
		return p.DeliveryDate != nil
	})).Return(nil)

	err := service.ProductADD(context.Background(), p)

	require.NoError(t, err)

	repo.AssertCalled(t, "ProductADD", mock.Anything, mock.Anything)
}
func TestProductUpdate_NotValid(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	p := models.ProductModel{}

	err := service.ProductUpdate(context.Background(), p)

	require.Error(t, err)

	repo.AssertNotCalled(t, "ProductUpdate", mock.Anything, mock.Anything)
}

func TestProductUpdate_StorageIDNil(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	p := models.ProductModel{
		Article:     1,
		ProductName: "test",
		Weight:      10,
		StorageID:   nil,
	}

	err := service.ProductUpdate(context.Background(), p)

	require.Error(t, err)
	require.Contains(t, err.Error(), "storage id")

	repo.AssertNotCalled(t, "ProductUpdate", mock.Anything, mock.Anything)
}

func TestProductUpdate_ExpireDateSet(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	storageID := 1
	now := time.Now()

	p := models.ProductModel{
		Article:     1,
		ProductName: "test",
		Weight:      10,
		StorageID:   &storageID,
		ExpireDate:  &now,
	}

	err := service.ProductUpdate(context.Background(), p)

	require.Error(t, err)
	require.Contains(t, err.Error(), "must be not expired")

	repo.AssertNotCalled(t, "ProductUpdate", mock.Anything, mock.Anything)
}

func TestProductUpdate_RepoError(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	storageID := 1

	p := models.ProductModel{
		Article:     1,
		ProductName: "test",
		Weight:      10,
		StorageID:   &storageID,
	}

	repo.On("ProductUpdate", mock.Anything, mock.AnythingOfType("models.ProductModel")).
		Return(errors.New("db error"))

	err := service.ProductUpdate(context.Background(), p)

	require.Error(t, err)
	require.EqualError(t, err, "db error")

	repo.AssertCalled(t, "ProductUpdate", mock.Anything, mock.AnythingOfType("models.ProductModel"))
}

func TestProductUpdate_Success(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	storageID := 1

	p := models.ProductModel{
		Article:     1,
		ProductName: "test",
		Weight:      10,
		StorageID:   &storageID,
		ExpireDate:  nil,
	}

	repo.On("ProductUpdate", mock.Anything, mock.AnythingOfType("models.ProductModel")).
		Return(nil)

	err := service.ProductUpdate(context.Background(), p)

	require.NoError(t, err)

	repo.AssertCalled(t, "ProductUpdate", mock.Anything, mock.AnythingOfType("models.ProductModel"))
}

func TestProductDelete_InvalidArticle(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	err := service.ProductDelete(context.Background(), 0)

	require.Error(t, err)
	require.EqualError(t, err, "article must be greater than zero")

	repo.AssertNotCalled(t, "ProductDelete", mock.Anything, mock.Anything)
}

func TestProductDelete_RepoError(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	repo.On("ProductDelete", mock.Anything, 1).
		Return(errors.New("db error"))

	err := service.ProductDelete(context.Background(), 1)

	require.Error(t, err)
	require.EqualError(t, err, "db error")

	repo.AssertCalled(t, "ProductDelete", mock.Anything, 1)
}

func TestProductDelete_Success(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	repo.On("ProductDelete", mock.Anything, 1).
		Return(nil)

	err := service.ProductDelete(context.Background(), 1)

	require.NoError(t, err)

	repo.AssertCalled(t, "ProductDelete", mock.Anything, 1)
}

func TestProductExpire_InvalidArticle(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	err := service.ProductExpire(context.Background(), 0)

	require.Error(t, err)
	require.EqualError(t, err, "article must be greater than zero")

	repo.AssertNotCalled(t, "GetProduct", mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "ProductExpire", mock.Anything, mock.Anything)
}

func TestProductExpire_NotFound(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	repo.On("GetProduct", mock.Anything, 1).
		Return(nil, sql.ErrNoRows)

	err := service.ProductExpire(context.Background(), 1)

	require.Error(t, err)
	require.EqualError(t, err, "product not found")

	repo.AssertCalled(t, "GetProduct", mock.Anything, 1)
	repo.AssertNotCalled(t, "ProductExpire", mock.Anything, mock.Anything)
}

func TestProductExpire_GetProductError(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	repo.On("GetProduct", mock.Anything, 1).
		Return(nil, errors.New("db error"))

	err := service.ProductExpire(context.Background(), 1)

	require.Error(t, err)
	require.EqualError(t, err, "db error")

	repo.AssertCalled(t, "GetProduct", mock.Anything, 1)
	repo.AssertNotCalled(t, "ProductExpire", mock.Anything, mock.Anything)
}

func TestProductExpire_AlreadyExpired(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	expire := time.Now()

	product := &models.ProductModel{
		Article:    1,
		ExpireDate: &expire,
	}

	repo.On("GetProduct", mock.Anything, 1).
		Return(product, nil)

	err := service.ProductExpire(context.Background(), 1)

	require.Error(t, err)
	require.EqualError(t, err, "product already expired")

	repo.AssertCalled(t, "GetProduct", mock.Anything, 1)
	repo.AssertNotCalled(t, "ProductExpire", mock.Anything, mock.Anything)
}

func TestProductExpire_RepoError(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	product := &models.ProductModel{
		Article: 1,
	}

	repo.On("GetProduct", mock.Anything, 1).
		Return(product, nil)

	repo.On("ProductExpire",
		mock.Anything,
		mock.MatchedBy(func(p models.ProductModel) bool {
			return p.Article == 1 && p.ExpireDate != nil
		}),
	).Return(errors.New("db error"))

	err := service.ProductExpire(context.Background(), 1)

	require.Error(t, err)
	require.EqualError(t, err, "db error")

	repo.AssertCalled(t, "GetProduct", mock.Anything, 1)
	repo.AssertCalled(t, "ProductExpire", mock.Anything, mock.Anything)
}

func TestProductExpire_Success(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	product := &models.ProductModel{
		Article: 1,
	}

	repo.On("GetProduct", mock.Anything, 1).
		Return(product, nil)

	repo.On("ProductExpire",
		mock.Anything,
		mock.MatchedBy(func(p models.ProductModel) bool {
			return p.Article == 1 && p.ExpireDate != nil
		}),
	).Return(nil)

	err := service.ProductExpire(context.Background(), 1)

	require.NoError(t, err)

	repo.AssertCalled(t, "GetProduct", mock.Anything, 1)
	repo.AssertCalled(t, "ProductExpire", mock.Anything, mock.Anything)
}

func TestGetProduct_InvalidArticle(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	product, err := service.GetProductFromProducts(context.Background(), 0)

	require.Error(t, err)
	require.EqualError(t, err, "id should be greater than zero")
	require.Nil(t, product)

	repo.AssertNotCalled(t, "GetProduct", mock.Anything, mock.Anything)
}

func TestGetProduct_NotFound(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	repo.On("GetProduct", mock.Anything, 1).
		Return(nil, sql.ErrNoRows)

	product, err := service.GetProductFromProducts(context.Background(), 1)

	require.Error(t, err)
	require.EqualError(t, err, "product not found")
	require.Nil(t, product)

	repo.AssertCalled(t, "GetProduct", mock.Anything, 1)
}

func TestGetProduct_RepoError(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	repo.On("GetProduct", mock.Anything, 1).
		Return(nil, errors.New("db error"))

	product, err := service.GetProductFromProducts(context.Background(), 1)

	require.Error(t, err)
	require.EqualError(t, err, "db error")
	require.Nil(t, product)

	repo.AssertCalled(t, "GetProduct", mock.Anything, 1)
}

func TestGetProduct_Success(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	expected := &models.ProductModel{
		Article:     1,
		ProductName: "test",
	}

	repo.On("GetProduct", mock.Anything, 1).
		Return(expected, nil)

	product, err := service.GetProductFromProducts(context.Background(), 1)

	require.NoError(t, err)
	require.Equal(t, expected, product)

	repo.AssertCalled(t, "GetProduct", mock.Anything, 1)
}

func TestProductsList_RepoError(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	repo.On("ListProducts", mock.Anything).
		Return(nil, errors.New("db error"))

	result, err := service.ProductsListFromProducts(context.Background())

	require.Error(t, err)
	require.EqualError(t, err, "db error")
	require.Nil(t, result)

	repo.AssertCalled(t, "ListProducts", mock.Anything)
}

func TestProductsList_Success(t *testing.T) {
	repo := new(mocks.MockWHouseRepository)
	service := NewProductService(repo, slog.Default())

	expected := []models.ProductModel{
		{Article: 1, ProductName: "test1"},
		{Article: 2, ProductName: "test2"},
	}

	repo.On("ListProducts", mock.Anything).
		Return(expected, nil)

	result, err := service.ProductsListFromProducts(context.Background())

	require.NoError(t, err)
	require.Equal(t, expected, result)

	repo.AssertCalled(t, "ListProducts", mock.Anything)
}
