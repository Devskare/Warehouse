package repository

import (
	"context"
	"fmt"
	"warehouse/modules/Whouse/models"

	"github.com/jmoiron/sqlx"
)

type WHouseRepositoryDB struct {
	db *sqlx.DB
}

func NewWHouseRepository(db *sqlx.DB) WHouser { return &WHouseRepositoryDB{db: db} }

func (r *WHouseRepositoryDB) ProductADD(ctx context.Context, product models.ProductModel) error {
	sqlQuery := `INSERT INTO products 
		(article, product_name, storage_id, delivery_date, expire_date, weight)
		VALUES (:article, :product_name, :storage_id, :delivery_date, :expire_date, :weight)`

	_, err := r.db.ExecContext(ctx, sqlQuery, product)
	if err != nil {
		return err
	}
	fmt.Println("Product Add Success")
	return nil
}

func (r *WHouseRepositoryDB) ProductUpdate(ctx context.Context, product models.ProductModel) error {
	sqlQuery := `UPDATE products
				 SET   product_name = $1, storage_id = $2, expire_date = $3, weight = $4
				 WHERE article = $5`

	_, err := r.db.ExecContext(ctx, sqlQuery, product)
	if err != nil {
		return err
	}
	fmt.Println("Product Update Success")
	return nil
}

func (r *WHouseRepositoryDB) ProductDelete(ctx context.Context, product models.ProductModel) error {
	sqlQuery := `DELETE FROM products WHERE article = $1`

	_, err := r.db.ExecContext(ctx, sqlQuery, product)
	if err != nil {
		return err
	}
	fmt.Println("Product Delete Success")
	return nil
}
