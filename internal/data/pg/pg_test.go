package pg

import (
	"testing"

	"github.com/dead-letter/dead-letter-data/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/goosemigrator"
	"github.com/stretchr/testify/assert"
)

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dbconf := pgtestdb.Config{
		DriverName: "pgx",
		User:       "postgres",
		Password:   "password",
		Host:       "localhost",
		Port:       "5433",
		Options:    "sslmode=disable",
	}
	m := goosemigrator.New(
		".",
		goosemigrator.WithFS(migrations.Files),
	)
	c := pgtestdb.Custom(t, dbconf, m)
	assert.NotEqual(t, dbconf, *c)

	pool, err := OpenPool(c.URL())
	assert.NoError(t, err)

	return pool
}
