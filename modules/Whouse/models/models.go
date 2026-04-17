package models

import (
	"errors"
	"time"
)

type ProductModel struct {
	ID           int        `db:"id" `
	Article      int        `db:"article" `
	ProductName  string     `db:"product_name" `
	StorageID    *int       `db:"storage_id" `
	DeliveryDate *time.Time `db:"delivery_date" `
	ExpireDate   *time.Time `db:"expire_date" `
	Weight       float64    `db:"weight" `
}

type StorageModel struct {
	ID            int     `db:"id" `
	MaxWeight     float64 `db:"max_weight" `
	CurrentWeight float64 `db:"current_weight" `
}

func (s *ProductModel) Validate() error {
	if s.Article < 0 {
		return errors.New("article must be greater than zero, invalid article number")
	}
	if len(s.ProductName) > 225 {
		return errors.New("product name must be less than 225 characters, invalid product name")
	}
	if len(s.ProductName) < 1 {
		return errors.New("product name must be greater than zero, invalid product name")
	}
	if s.Weight < 0 {
		return errors.New("weight must be greater than zero, invalid product weight")
	}
	return nil
}
