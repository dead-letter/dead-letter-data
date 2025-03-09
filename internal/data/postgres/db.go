package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/dead-letter/dead-letter-data/internal/data"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	*data.DB
	Pool *pgxpool.Pool
}

// Create a database pool. Don't forget to Close()
func NewPostgresDB(dsn string) (*PostgresDB, error) {
	pool, err := openPool(dsn)
	if err != nil {
		return nil, err
	}

	pg := &PostgresDB{
		DB: &data.DB{
			Users:   &UserRepository{pool},
			Riders:  &RiderRepository{pool},
			Vendors: &VendorRepository{pool},
		},
		Pool: pool,
	}

	return pg, nil
}

// Close closes all connections in the database pool
func (pg *PostgresDB) Close() {
	pg.Pool.Close()
}

func openPool(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func pgErrCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}

	return ""
}
