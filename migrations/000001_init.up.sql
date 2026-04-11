CREATE TABLE IF NOT EXISTS storages (
                                        id SERIAL PRIMARY KEY,
                                        max_weight NUMERIC(15, 3) NOT NULL CHECK (max_weight >= 0),
    current_weight NUMERIC(15,3) NOT NULL DEFAULT 0 CHECK (current_weight >= 0)
    );


CREATE TABLE IF NOT EXISTS products (
                                        id SERIAL PRIMARY KEY,
                                        article INTEGER NOT NULL UNIQUE,
                                        product_name VARCHAR(255) NOT NULL,
    storage_id INTEGER,
    delivery_date DATE,
    expire_date DATE,
    weight NUMERIC(15,3) NOT NULL CHECK (weight > 0),

    CONSTRAINT fk_products_storage
    FOREIGN KEY (storage_id)
    REFERENCES storages(id)
    ON DELETE RESTRICT
    );