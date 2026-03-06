package models

import "time"

type ProductModel struct {
	ID           int
	Article      int
	ProductName  string
	StorageID    int
	DeliveryDate time.Time
	ExpireDate   time.Time
	Weight       float64
}
