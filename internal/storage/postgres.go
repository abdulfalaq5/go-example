package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/falaqmsi/go-example/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB holds the two named connection pools used across the application.
type DB struct {
	Main  *pgxpool.Pool // general-purpose / transactional data
	Audit *pgxpool.Pool // append-only audit / event log data
}

// Connect opens and validates both PostgreSQL connection pools.
// Caller is responsible for calling DB.Close() when done.
func Connect(ctx context.Context, cfg config.DBConfig) (*DB, error) {
	main, err := newPool(ctx, cfg, cfg.MainDSN, "DBMain")
	if err != nil {
		return nil, err
	}

	audit, err := newPool(ctx, cfg, cfg.AuditDSN, "DBAudit")
	if err != nil {
		main.Close()
		return nil, err
	}

	return &DB{Main: main, Audit: audit}, nil
}

// Close gracefully shuts down both pools. Safe to call on a nil receiver.
func (db *DB) Close() {
	if db == nil {
		return
	}
	if db.Main != nil {
		db.Main.Close()
		log.Println("[storage] DBMain pool closed")
	}
	if db.Audit != nil {
		db.Audit.Close()
		log.Println("[storage] DBAudit pool closed")
	}
}

// newPool builds and pings a pgxpool using the shared pool settings.
func newPool(ctx context.Context, cfg config.DBConfig, dsn, name string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("[storage] %s: parse DSN: %w", name, err)
	}

	// ── Pool tuning ────────────────────────────────────────────────────────────
	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = cfg.HealthCheckPeriod

	// Timeout for acquiring a connection from the pool
	poolCfg.MaxConnLifetimeJitter = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("[storage] %s: create pool: %w", name, err)
	}

	// Verify connectivity immediately so startup fails fast if DB is unreachable.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("[storage] %s: ping failed: %w", name, err)
	}

	log.Printf("[storage] %s connected (maxConns=%d minConns=%d)",
		name, cfg.MaxConns, cfg.MinConns)

	return pool, nil
}
