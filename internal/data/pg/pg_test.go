package pg

import (
	"context"
	"testing"
	"time"

	"github.com/dead-letter/dead-letter-data/migrations"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var s UserService

func TestMain(m *testing.M) {
	// Start Postgres instance
	var err error
	ctx := context.Background()
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(1*time.Minute)),
	)
	if err != nil {
		panic(err)
	}
	defer pgContainer.Terminate(ctx)

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	pool, err := OpenPool(dsn)
	if err != nil {
		panic(err)
	}

	// Run migrations
	db := stdlib.OpenDBFromPool(pool)
	migrations.Up(db)
	db.Close()

	s = UserService{Pool: pool}

	m.Run()
}
