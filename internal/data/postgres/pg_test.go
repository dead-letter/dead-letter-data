package postgres

import (
	"testing"

	"github.com/dead-letter/dead-letter-data/migrations"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/goosemigrator"
	"github.com/stretchr/testify/assert"
)

func testDB(t *testing.T) *PostgresDB {
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

	pg, err := NewPostgresDB(c.URL())
	assert.NoError(t, err)

	return pg
}
