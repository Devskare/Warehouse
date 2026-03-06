package connection

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// "postgres://YourUserName:YourPassword@YourHostName:5432/YourDatabaseName"

func Connection(ctx context.Context) (*pgx.Conn, error) {
	defer fmt.Println("Successfully connected to PostgreSQL")

	return pgx.Connect(ctx, "postgres://postgres:1234@localhost:5432/postgres")

}
