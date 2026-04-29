package repository

import (
	"context"

	"github.com/Devskare/Warehouse/modules/Whouse/models"
)

type WHouser interface {
	ProductADD(ctx context.Context, product models.ProductModel) error
	ProductUpdate(ctx context.Context, product models.ProductModel) error
	ProductDelete(ctx context.Context, article int) error
	StorageADD(ctx context.Context, MaxWeight float64) error
	ListProducts(ctx context.Context) ([]models.ProductModel, error)
	ListStorages(ctx context.Context) ([]models.StorageModel, error)
	GetProduct(ctx context.Context, article int) (*models.ProductModel, error)
	ProductExpire(ctx context.Context, product models.ProductModel) error
}
