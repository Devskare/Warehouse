package connection

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

// "postgres://YourUserName:YourPassword@YourHostName:5432/YourDatabaseName"

func Connection(ctx context.Context) (*pgx.Conn, error) {
	ConnString := os.Getenv("CONNECTION_STRING")

	return pgx.Connect(ctx, ConnString)

}
