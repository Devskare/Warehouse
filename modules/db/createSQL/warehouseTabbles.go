package createSQL

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateTables(ctx context.Context, conn *pgx.Conn) error {
	sqlStorage := `
	CREATE TABLE IF NOT EXISTS storages (
	    id SERIAL PRIMARY KEY,
	  max_weight FLOAT
	);
	`
	sqlProducts := `
CREATE TABLE IF NOT EXISTS products (
id SERIAL PRIMARY KEY,
article INTEGER NOT NULL UNIQUE,
product_name VARCHAR(255) NOT NULL,
storage_id INTEGER,
delivery_date DATE,
expire_date DATE,
weight FLOAT NOT NULL
    );
	`

	_, err := conn.Exec(ctx, sqlStorage)
	if err != nil {

		return err
	}
	_, err = conn.Exec(ctx, sqlProducts)
	if err != nil {

		return err
	}
	return nil
}
