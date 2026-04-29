package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/Devskare/Warehouse/modules/Whouse/models"
	"github.com/Devskare/Warehouse/modules/Whouse/repository"
)

type ProductService struct {
	repo repository.WHouser
	log  *slog.Logger
}

func NewProductService(repo repository.WHouser, log *slog.Logger) *ProductService {
	return &ProductService{repo: repo, log: log}
}

func (s *ProductService) ProductADD(ctx context.Context, p models.ProductModel) error {
	err := p.Validate()
	if err != nil {
		s.log.Error("Product validation failed", "error", err.Error())
		return err
	}
	if p.StorageID == nil {
		s.log.Error("Product storage ID is nil")
		return errors.New("product storage ID is nil")

	}
	t := time.Now().UTC()
	p.DeliveryDate = &t
	err = s.repo.ProductADD(ctx, p)
	if err != nil {
		s.log.Error("Product add failed", "error", err.Error())
		return err
	}
	return nil
}

func (s *ProductService) ProductUpdate(ctx context.Context, p models.ProductModel) error {
	err := p.Validate()
	if err != nil {
		s.log.Error("Product validation failed", "error", err.Error())
		return err
	}
	if p.StorageID == nil || p.ExpireDate != nil {
		s.log.Error("Product storage ID is nil or product has been expired", "error", "Product storage ID is nil or expired")
		return errors.New("Product  should have storage id and must be not expired")
	}
	err = s.repo.ProductUpdate(ctx, p)
	if err != nil {
		s.log.Error("Product update failed", "error", err.Error())
		return err
	}
	return nil
}

func (s *ProductService) ProductDelete(ctx context.Context, article int) error {
	if article < 1 {
		s.log.Error("Product article must be greater than zero")
		return errors.New("article must be greater than zero")
	}
	err := s.repo.ProductDelete(ctx, article)
	if err != nil {
		s.log.Error("Product delete failed", "error", err.Error())
		return err
	}
	return nil

}

func (s *ProductService) ProductExpire(ctx context.Context, article int) error {

	if article < 1 {
		s.log.Error("Product article must be greater than zero")
		return errors.New("article must be greater than zero")
	}

	product, err := s.repo.GetProduct(ctx, article)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Error("Product not found", "error", err.Error())
			return errors.New("product not found")
		}

		s.log.Error("Failed to get product", "error", err.Error())
		return err
	}

	if product.ExpireDate != nil {
		s.log.Error("Product already expired")
		return errors.New("product already expired")
	}

	now := time.Now().UTC()
	product.ExpireDate = &now

	err = s.repo.ProductExpire(ctx, *product)
	if err != nil {
		s.log.Error("Product expire failed", "error", err.Error())
		return err
	}

	return nil
}

func (s *ProductService) GetProductFromProducts(ctx context.Context, article int) (*models.ProductModel, error) {
	if article < 1 {
		s.log.Error("Product id should be greater than zero")
		return nil, errors.New("id should be greater than zero")
	}
	product, err := s.repo.GetProduct(ctx, article)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Error("Product not found", "error", err.Error())
			return nil, errors.New("product not found")
		}
		s.log.Error("Product get failed", "error", err.Error())
		return nil, err
	}
	return product, nil
}

func (s *ProductService) ProductsListFromProducts(ctx context.Context) ([]models.ProductModel, error) {
	productList, err := s.repo.ListProducts(ctx)
	if err != nil {
		s.log.Error("Products list failed", "error", err.Error())
		return nil, err
	}
	return productList, nil
}
