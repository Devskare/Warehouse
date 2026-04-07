package models

import "time"

type ProductModel struct {
	ID           int        `db:"id" `
	Article      int        `db:"article" `
	ProductName  string     `db:"product_name" `
	StorageID    int        `db:"storage_id" `
	DeliveryDate *time.Time `db:"delivery_date" `
	ExpireDate   *time.Time `db:"expire_date" `
	Weight       float64    `db:"weight" `
}

type StorageModel struct {
	ID            int     `db:"id" `
	MaxWeight     float64 `db:"max_weight" `
	CurrentWeight float64 `db:"current_weight" `
}
