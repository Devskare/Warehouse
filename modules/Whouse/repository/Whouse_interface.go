package repository

import (
	"context"
	"warehouse/modules/Whouse/models"
)

type WHouser interface {
	ProductADD(ctx context.Context, product models.ProductModel) error
	ProductUpdate(ctx context.Context, product models.ProductModel) error
	ProductDelete(ctx context.Context, ID int) error
}
