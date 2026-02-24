package connection

import (
	context2 "context"
	"cust-rtmn-orch-library/resources/properties"
	"cust-rtmn-orch-library/resources/secrets"
	"fmt"
	"sync"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
)

var connections = make(map[string]*connection)
var mutex sync.RWMutex

type connection struct {
	Pool *pgxpool.Pool
}

func getPool(ctx context2.Context, database string) (*pgxpool.Pool, error) {
	ip := secrets.GetSecretValue("pgIp")
	if properties.GetProperty().IsLocal() {
		ip = "localhost"
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		ip,
		secrets.GetSecretValue("pgPort"),
		secrets.GetSecretValue("pgUser"),
		secrets.GetSecretValue("pgPass"),
		database)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	config.MaxConnLifetime = time.Minute * 1
	config.MaxConnLifetimeJitter = time.Second * 5
	config.MaxConnIdleTime = time.Minute * 1
	config.HealthCheckPeriod = time.Second * 30
	config.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTracerProvider(otel.GetTracerProvider()))
	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func PerformConnection(ctx context2.Context, database string) (pool *pgxpool.Pool, err error) {
	mutex.Lock()
	connection := connections[database]
	mutex.Unlock()
	if connection != nil && connection.Pool != nil {
		pool = connection.Pool
	}
	if pool != nil {
		return pool, nil
	}
	pool, err = getPool(ctx, database)
	if err == nil {
		mutex.Lock()
		connections[database] = cacheConnection(pool)
		mutex.Unlock()
	}
	return pool, err
}

func cacheConnection(pool *pgxpool.Pool) *connection {
	return &connection{Pool: pool}
}
