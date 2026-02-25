package connection

import (
	context2 "context"
	"database/sql"
	"fmt"
	"log"
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

func ConnectDb(username string, password string, host string, port string, dbname string) (*sql.DB, error) {

	// Build connection string with Sprintf
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		username, password, host, port, dbname,
	)

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}
	defer db.Close() // Always close when done

	// Verify the connection is alive
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return db, nil
}

func getPool(ctx context2.Context, database string, pgUser string, pgPassword string, pgIp string, port string) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgIp,
		port,
		pgUser,
		pgPassword,
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

func PerformConnection(ctx context2.Context, database string, pgUser string, pgPassword string, pgIp string, port string) (pool *pgxpool.Pool, err error) {
	mutex.Lock()
	connection := connections[database]
	mutex.Unlock()
	if connection != nil && connection.Pool != nil {
		pool = connection.Pool
	}
	if pool != nil {
		return pool, nil
	}
	pool, err = getPool(ctx, database, pgUser, pgPassword, pgIp, port)
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
