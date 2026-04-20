package service

import (
	"context"
	"warehouse/modules/Whouse/models"
)

type Product interface {
	ProductADD(ctx context.Context, p models.ProductModel) error
	ProductDEL(ctx context.Context, p models.ProductModel) error
	ProductEXPIRE(ctx context.Context, p models.ProductModel) error
	ProductUPDATE(ctx context.Context, p models.ProductModel) error
	GetProductFromProducts(ctx context.Context, article int) (*models.ProductModel, error)
	ProductsListFromProducts(ctx context.Context) ([]models.ProductModel, error)
	ProductExpire(ctx context.Context, article int) error
}

type Storage interface {
	StorageADD(ctx context.Context, MaxWeight float64) error
	ListStorages(ctx context.Context) ([]models.StorageModel, error)
}
