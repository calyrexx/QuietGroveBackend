package postgresdb

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Calyr3x/QuietGrooveBackend/internal/configuration"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"math"
	"os"
	"strings"
)

const (
	filePath = "deploy/postgresddl.sql"
)

func NewPostgres(ctx context.Context, c configuration.Postgres) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		c.User, c.Password, c.Host, c.Port, c.Database,
	))
	if err != nil {
		return nil, err
	}

	if c.MinConnections < math.MinInt32 || c.MinConnections > math.MaxInt32 {
		return nil, fmt.Errorf("MinConnections value %d is out of range for int32", c.MinConnections)
	}

	if c.MaxConnections < math.MinInt32 || c.MaxConnections > math.MaxInt32 {
		return nil, fmt.Errorf("MaxConnections value %d is out of range for int32", c.MaxConnections)
	}

	// nolint:gosec // выше проверена возможность переполнения
	config.MinConns = int32(c.MinConnections)
	config.MaxConns = int32(c.MaxConnections)
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

func InitDbInst(ctx context.Context, c *pgxpool.Pool) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	// Чтение и выполнение запросов построчно
	scanner := bufio.NewScanner(file)
	var queryBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		queryBuilder.WriteString(line + "\n")
		if strings.HasSuffix(line, ";") {
			query := queryBuilder.String()
			_, err := c.Exec(ctx, query)
			if err != nil {
				log.Fatalf("Error executing query: %v\nQuery: %s", err, query)
			}
			queryBuilder.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
