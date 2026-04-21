package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	"warehouse/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// "postgres://YourUserName:YourPassword@YourHostName:5432/YourDatabaseName"

func NewSqlDB(logger *slog.Logger, cfg *config.DB) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
	var dbRow *sql.DB
	var err error

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(time.Second * time.Duration(cfg.Timeout))

	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("timeout %d exceeded", cfg.Timeout)

		case <-ticker.C:
			dbRow, err = sql.Open(cfg.Driver, dsn)
			if err != nil {
				logger.Error("failed to open sql data base", slog.Any("error", err))
				return nil, err
			}
			err = dbRow.Ping()
			if err == nil {
				db := sqlx.NewDb(dbRow, cfg.Driver)
				db.SetMaxOpenConns(cfg.MaxConn)
				return db, nil
			}
			logger.Error("failed to connect sql data base", slog.Any("error", err))

		}

	}
}
