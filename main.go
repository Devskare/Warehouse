package main

import (
	"context"
	"time"
	"warehouse/models"
	"warehouse/postgres/connection"
	_ "warehouse/postgres/connection"
	"warehouse/postgres/createSQL"
	"warehouse/service"
)

func main() {

	ctx := context.Background()
	conn, err := connection.Connection(ctx)
	if err != nil {
		panic(err)
	}
	if err := createSQL.CreateTables(ctx, conn); err != nil {
		panic(err)

	}
	var prod models.ProductModel

	prod.Article = 123
	prod.ProductName = "abc"
	prod.StorageID = 1
	prod.DeliveryDate = time.Now()
	prod.ExpireDate = time.Now().AddDate(0, 0, 1)
	prod.Weight = 1
	err = service.ProductADD(ctx, conn, prod)
	if err != nil {
		panic(err)
	}

	err = service.ProductDelete(ctx, conn, prod)
	if err != nil {
		panic(err)
	}

	err = service.ProductADD(ctx, conn, prod)
	if err != nil {
		panic(err)
	}
	prod.ProductName = "cba"
	prod.StorageID = 2
	prod.ExpireDate = time.Now().AddDate(1, 0, 0)
	prod.Weight = 5
	err = service.ProductUpdate(ctx, conn, prod)
	if err != nil {
		panic(err)
	}
}
