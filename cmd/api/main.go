package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/dead-letter/dead-letter-data/internal/data/postgres"
	"github.com/dead-letter/dead-letter-data/internal/grpc"
	"github.com/dead-letter/dead-letter-data/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/lmittmann/tint"
)

type config struct {
	dev  bool
	port int
	dsn  string
}

func main() {
	var cfg config

	// Default for prod
	flag.BoolVar(&cfg.dev, "dev", false, "Development mode")
	flag.IntVar(&cfg.port, "port", 50051, "API server port")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")
	flag.Parse()

	h := newSlogHandler(cfg.dev)
	logger := slog.New(h)

	pg, err := postgres.NewDB(cfg.dsn)
	if err != nil {
		fatal(logger, err)
	}
	defer pg.Close()

	err = runMigrations(pg.Pool, cfg.dev)
	if err != nil {
		fatal(logger, err)
	}

	addr := fmt.Sprintf(":%d", cfg.port)
	srv := grpc.NewServer(addr, pg.DB)
	err = srv.ListenAndServe()
	if err != nil {
		fatal(logger, err)
	}
}

func newSlogHandler(dev bool) slog.Handler {
	if dev {
		// Development text hanlder
		return tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		})
	}

	// Production use JSON handler with default opts
	return slog.NewJSONHandler(os.Stdout, nil)
}

func runMigrations(pool *pgxpool.Pool, reset bool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	err := migrations.Up(db)
	if err != nil {
		return err
	}

	if reset {
		err = migrations.Reset(db)
		if err != nil {
			return err
		}
	}

	return nil
}

func fatal(logger *slog.Logger, err error) {
	logger.Error("fatal", slog.Any("err", err))
	os.Exit(1)
}
