package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

const defaultTimeout = 5 * time.Second

// NewDatabasePool initializes a new database connection pool.
func NewDatabasePool(ctx context.Context, dbUrl string, logger *slog.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		logger.Error("failed to parse database URL", err)
		return nil, err
	}

	// Set important pool configurations close to defaults
	config.MaxConns = 10                      // Maximum number of connections in the pool
	config.MinConns = 1                       // Minimum number of connections in the pool
	config.MaxConnLifetime = time.Hour        // Maximum lifetime of a connection
	config.MaxConnIdleTime = 30 * time.Minute // Maximum idle time for a connection
	config.HealthCheckPeriod = time.Minute    // How often to check the health of idle connections

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error("failed to connect to database", err)
		return nil, err
	}

	logger.Info("connected to database")
	return pool, nil
}
