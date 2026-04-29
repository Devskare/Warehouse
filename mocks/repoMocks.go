package mocks

import (
	"context"

	"github.com/Devskare/Warehouse/modules/Whouse/models"
	"github.com/Devskare/Warehouse/modules/Whouse/repository"

	"github.com/stretchr/testify/mock"
)

type MockWHouseRepository struct {
	mock.Mock
}

func (m *MockWHouseRepository) ProductADD(ctx context.Context, product models.ProductModel) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockWHouseRepository) ProductUpdate(ctx context.Context, p models.ProductModel) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockWHouseRepository) ProductDelete(ctx context.Context, article int) error {
	args := m.Called(ctx, article)
	return args.Error(0)
}

func (m *MockWHouseRepository) GetProduct(ctx context.Context, article int) (*models.ProductModel, error) {
	args := m.Called(ctx, article)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductModel), args.Error(1)
}

func (m *MockWHouseRepository) ListProducts(ctx context.Context) ([]models.ProductModel, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.ProductModel), args.Error(1)
}

func (m *MockWHouseRepository) StorageADD(ctx context.Context, maxWeight float64) error {
	args := m.Called(ctx, maxWeight)
	return args.Error(0)
}

func (m *MockWHouseRepository) ListStorages(ctx context.Context) ([]models.StorageModel, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.StorageModel), args.Error(1)
}

func (m *MockWHouseRepository) ProductExpire(ctx context.Context, product models.ProductModel) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

var _ repository.WHouser = (*MockWHouseRepository)(nil)
