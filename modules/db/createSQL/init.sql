CREATE TABLE IF NOT EXISTS storages (
                                        id SERIAL PRIMARY KEY,
                                        max_weight FLOAT
);


CREATE TABLE IF NOT EXISTS products (
                                        id SERIAL PRIMARY KEY,
                                        article INTEGER NOT NULL UNIQUE,
                                        product_name VARCHAR(255) NOT NULL,
    storage_id INTEGER,
    delivery_date DATE,
    expire_date DATE,
    weight FLOAT NOT NULL
    );