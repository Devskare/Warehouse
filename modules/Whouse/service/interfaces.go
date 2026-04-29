package service

import (
	"context"

	"github.com/devskar/warehouse/modules/Whouse/models"
)

type Product interface {
	ProductADD(ctx context.Context, p models.ProductModel) error
	ProductUpdate(ctx context.Context, p models.ProductModel) error
	ProductDelete(ctx context.Context, article int) error
	ProductExpire(ctx context.Context, article int) error
	GetProductFromProducts(ctx context.Context, article int) (*models.ProductModel, error)
	ProductsListFromProducts(ctx context.Context) ([]models.ProductModel, error)
}

type Storage interface {
	StorageADD(ctx context.Context, MaxWeight float64) error
	ListStorages(ctx context.Context) ([]models.StorageModel, error)
}
