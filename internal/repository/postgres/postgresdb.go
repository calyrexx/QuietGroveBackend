package postgresdb

import (
	"context"
	"fmt"
	"github.com/Calyr3x/QuietGrooveBackend/internal/configuration"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"os"
)

const (
	filePath = "deploy/postgres.sql"
)

func NewPostgres(ctx context.Context, c configuration.Postgres) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		c.User, c.Password, c.Host, c.Port, c.Database,
	))
	if err != nil {
		return nil, err
	}

	config.MinConns = c.MinConnections
	config.MaxConns = c.MaxConnections
	config.MaxConnIdleTime = c.IdleConnection
	config.MaxConnLifetime = c.LifeTimeConnection
	config.MaxConnLifetimeJitter = c.JitterConnection

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	err = InitDbInst(ctx, pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func InitDbInst(ctx context.Context, pool *pgxpool.Pool) error {

	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	mrr := conn.Conn().PgConn().Exec(ctx, string(sqlBytes))
	_, err = mrr.ReadAll()

	return err
}
