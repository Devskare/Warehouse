package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Devskare/Warehouse/modules/Whouse/models"

	"github.com/jmoiron/sqlx"
)

/* вывод - временное решение, пока не появятся сервисы. После этого будем логгировать. */
type WHouseRepositoryDB struct {
	db *sqlx.DB
}

func NewWHouseRepository(db *sqlx.DB) WHouser { return &WHouseRepositoryDB{db: db} }

func (r *WHouseRepositoryDB) StorageADD(ctx context.Context, MaxWeight float64) error {
	sqlQuery := `INSERT INTO storages
				 (max_weight)
			      VALUES ($1)`

	_, err := r.db.ExecContext(ctx, sqlQuery, MaxWeight)
	if err != nil {
		return err
	}
	fmt.Println("Added new storage")
	return nil
}

func (r *WHouseRepositoryDB) ProductADD(ctx context.Context, product models.ProductModel) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	sqlQuery := `UPDATE storages SET current_weight = current_weight + $1
				 WHERE id = $2
		  		 AND current_weight + $1 <= max_weight
				 RETURNING current_weight`

	var newWeight float64
	err = tx.GetContext(ctx, &newWeight, sqlQuery, product.Weight, *product.StorageID)
	if err != nil {
		return fmt.Errorf("storage overflow or not found: %w", err)
	}

	sqlQuery = `INSERT INTO products 
				(article, product_name, storage_id, delivery_date, expire_date, weight)
				VALUES (:article, :product_name, :storage_id, :delivery_date, :expire_date, :weight)`

	_, err = tx.NamedExecContext(ctx, sqlQuery, product)
	if err != nil {
		return err
	}
	fmt.Println("Added new product")
	return tx.Commit()

}

func (r *WHouseRepositoryDB) ProductUpdate(ctx context.Context, product models.ProductModel) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	sqlQuery := `SELECT weight, storage_id 
				 FROM products 
				 WHERE article = $1`

	var oldWeight float64
	var oldStorageID int
	err = tx.QueryRowContext(ctx, sqlQuery, product.Article).Scan(&oldWeight, &oldStorageID)
	if err != nil {
		return fmt.Errorf("cant found new old weight or old storageID in storage.  %w", err)
	}
	sqlQuery = `UPDATE storages
				SET current_weight = current_weight - $1
				WHERE id = $2
				RETURNING current_weight`
	var garbage float64
	err = tx.GetContext(ctx, &garbage, sqlQuery, oldWeight, oldStorageID)
	if err != nil {
		return fmt.Errorf("failed to update old storage: %w", err)
	}

	sqlQuery = `UPDATE storages
				SET current_weight = current_weight + $1
				WHERE id = $2
		  		AND current_weight + $1 <= max_weight
				RETURNING current_weight;`

	err = tx.GetContext(ctx, &garbage, sqlQuery, product.Weight, *product.StorageID)
	if err != nil {
		return fmt.Errorf("new storage is overflow!: %w", err)
	}

	sqlQuery = `UPDATE products
				SET product_name=:product_name, storage_id=:storage_id, delivery_date=:delivery_date, expire_date=:expire_date, weight=:weight
				WHERE article=:article`

	_, err = tx.NamedExecContext(ctx, sqlQuery, product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	fmt.Println("Updated product")
	return tx.Commit()

}

func (r *WHouseRepositoryDB) ProductExpire(ctx context.Context, product models.ProductModel) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var weight float64
	var storageID int

	err = tx.QueryRowContext(ctx, `
		SELECT weight, storage_id
		FROM products
		WHERE article = $1
		 AND expire_date IS NULL
		FOR UPDATE
	`, product.Article).Scan(&weight, &storageID)
	if err != nil {
		return fmt.Errorf("failed to lock product: %w", err)
	}

	sqlQuery := `UPDATE storages
				 SET current_weight = current_weight - $1
			 	 WHERE id = $2`

	_, err = tx.ExecContext(ctx, sqlQuery, weight, storageID)
	if err != nil {
		return fmt.Errorf("failed to update storage weight: %w", err)
	}

	sqlQuery = `UPDATE products
				SET storage_id = NULL,
				    expire_date = $1
				WHERE article = $2`

	_, err = tx.ExecContext(ctx, sqlQuery, product.ExpireDate, product.Article)
	if err != nil {
		return fmt.Errorf("failed to update product storage_id: %w", err)
	}

	return tx.Commit()

}

func (r *WHouseRepositoryDB) ProductDelete(ctx context.Context, article int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var weight float64
	var storageID int
	sqlQuery := `SELECT weight, storage_id 
				 FROM products 
				 WHERE article = $1`

	err = tx.QueryRowxContext(ctx, sqlQuery, article).Scan(&weight, &storageID)
	if err != nil {
		return fmt.Errorf("cant found weight and storageID: %w", err)
	}

	sqlQuery = `UPDATE storages
				SET current_weight = current_weight - $1
				WHERE id = $2
				RETURNING current_weight`
	var garbage float64
	err = tx.GetContext(ctx, &garbage, sqlQuery, weight, storageID)
	if err != nil {
		return fmt.Errorf("failed to delete old weight: %w", err)
	}

	sqlQuery = `DELETE FROM products 
      			WHERE article = $1`
	_, err = tx.ExecContext(ctx, sqlQuery, article)
	if err != nil {
		return err
	}

	fmt.Println("Deleted product")
	return tx.Commit()
}

func (r *WHouseRepositoryDB) ListProducts(ctx context.Context) ([]models.ProductModel, error) {
	var AllProducts []models.ProductModel
	sqlQuery := `SELECT id, article, product_name, storage_id, delivery_date, expire_date, weight 
				 FROM products
				 ORDER BY id`

	err := r.db.SelectContext(ctx, &AllProducts, sqlQuery)
	if err != nil {
		return nil, err
	}
	return AllProducts, nil
}

func (r *WHouseRepositoryDB) ListStorages(ctx context.Context) ([]models.StorageModel, error) {
	var AllStorages []models.StorageModel
	sqlQuery := `SELECT id, max_weight, current_weight 
				 FROM storages
				 ORDER BY id`

	err := r.db.SelectContext(ctx, &AllStorages, sqlQuery)
	if err != nil {
		return nil, err
	}
	return AllStorages, nil
}

func (r *WHouseRepositoryDB) GetProduct(ctx context.Context, article int) (*models.ProductModel, error) {
	var product models.ProductModel
	sqlQuery := `SELECT id, article, product_name, storage_id, delivery_date, expire_date, weight 
				 FROM products 
				 WHERE article = $1`

	err := r.db.GetContext(ctx, &product, sqlQuery, article)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return &product, nil
}
