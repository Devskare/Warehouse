package service

import (
	"context"
	"fmt"
	"warehouse/models"

	"github.com/jackc/pgx/v5"
)

func ProductADD(ctx context.Context, conn *pgx.Conn, product models.ProductModel) error {
	defer fmt.Println("Successfully add to PostgreSQL")
	sqlQiery := `INSERT INTO products ( article, product_name, storage_id , delivery_date, expire_date, weight)
    Values ( $1, $2, $3, $4, $5, $6) `

	_, err := conn.Exec(
		ctx, sqlQiery, product.Article, product.ProductName, product.StorageID, product.DeliveryDate, product.ExpireDate, product.Weight)
	if err != nil {
		return err
	}
	return nil
}

func ProductDelete(ctx context.Context, conn *pgx.Conn, product models.ProductModel) error {

	defer fmt.Println("Successfully delete from PostgreSQL")
	sqlQiery := `DELETE FROM products WHERE article = $1`

	_, err := conn.Exec(ctx, sqlQiery, product.Article)
	if err != nil {
		return err
	}
	return nil
}

func ProductUpdate(ctx context.Context, conn *pgx.Conn, product models.ProductModel) error {

	defer fmt.Println("Successfully update from PostgreSQL")
	sqlQiert := `UPDATE products 
				 SET   product_name = $1, storage_id = $2, expire_date = $3, weight = $4
				 WHERE article = $5`

	_, err := conn.Exec(ctx, sqlQiert, product.ProductName, product.StorageID, product.ExpireDate, product.Weight, product.Article)
	if err != nil {
		return err
	}
	return nil
}
