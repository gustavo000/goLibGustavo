package query_builder

import (
	context2 "context"
	"encoding/json"
	"log"
	"os"

	"github.com/gustavo000/goLibGustavo/pkg/functions"
	"github.com/gustavo000/goLibGustavo/pkg/pg-repository/connection"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (q *query) setOnContext(key string, value any) {}

func (q *query) generateMap(columns []Column) map[string]any {
	mapped := make(map[string]any)

	var waitGroupMap sync.WaitGroup
	waitGroupMap.Add(len(columns))
	mutex := &sync.RWMutex{}
	for _, column := range columns {
		chanValue := make(chan any)
		go func(column Column, chanValue chan any) {
			defer waitGroupMap.Done()
			var value = column.Value
			if q.queryType == INSERT && column.Generate {
				switch column.Type {
				case "uuid":
					value = uuid.NewString()
				case "date":
					value = time.Now().UTC()
				}
			}
			switch column.ParseTo {
			case "string":
				marshal, err := json.Marshal(value)
				if err == nil {
					chanValue <- string(marshal)
				}
			default:
				chanValue <- value
			}

		}(column, chanValue)
		value := <-chanValue
		if value != nil && !column.Serial {
			mutex.Lock()
			mapped[column.Name] = value
			mutex.Unlock()
		}
	}
	waitGroupMap.Wait()
	return mapped
}

func (q *query) generateColumns(object any) []Column {
	var mapString map[string]any
	var columns []Column
	err := functions.ParseTo(object, &mapString)
	if err == nil {
		typeOf := reflect.TypeOf(object)
		var mutex sync.RWMutex
		for i := 0; i < typeOf.NumField(); i++ {
			reflectField := typeOf.Field(i)
			column := GetColumn(reflectField, mapString)
			mutex.Lock()
			columns = append(columns, column)
			mutex.Unlock()
		}
	}
	return columns
}

func getDbConfig() (user, password, host, port, name string) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	user = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	name = os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	return
}

func (q *query) executeQuery() (rows pgx.Rows, conn *pgxpool.Conn, err error) {
	var spanCtx = q.ctx
	if spanCtx == nil {
		spanCtx = context2.Background()
	}
	resultQuery := strings.Join(q.queryParts, " ") + ";"
	u, p, h, prt, name := getDbConfig()
	pool, err := connection.PerformConnection(spanCtx, name, u, p, h, prt)
	if err != nil {
		return rows, conn, err
	}
	conn, err = pool.Acquire(spanCtx)
	if err != nil {
		return rows, conn, err
	}
	rows, err = conn.Query(spanCtx, resultQuery, q.namedArgs)
	if err != nil {
		return rows, conn, err
	}
	return rows, conn, err
}

func (q *query) executeCommand() (commandTag pgconn.CommandTag, conn *pgxpool.Conn, err error) {
	var spanCtx = q.ctx
	if spanCtx == nil {
		spanCtx = context2.Background()
	}
	resultQuery := strings.Join(q.queryParts, " ") + ";"
	u, p, h, prt := getDbConfig()
	pool, err := connection.PerformConnection(spanCtx, q.database, u, p, h, prt)
	if err != nil {
		return commandTag, conn, err
	}
	conn, err = pool.Acquire(spanCtx)
	if err != nil {
		return commandTag, conn, err
	}
	commandTag, err = conn.Exec(spanCtx, resultQuery, q.namedArgs)
	if err != nil {
		return commandTag, conn, err
	}
	return commandTag, conn, nil
}

func toChar(i int) rune {
	return rune('A' - 1 + i)
}

func (q *query) getResultFromRow(rows pgx.Rows) (any, error) {
	values, err := rows.Values()
	if err != nil {
		return nil, err
	}
	mutex := &sync.RWMutex{}
	group := &sync.WaitGroup{}
	typeOf := reflect.TypeOf(q.model)
	pointer := reflect.New(typeOf)
	group.Add(typeOf.NumField())
	for i := 0; i < typeOf.NumField(); i++ {
		go func(i int, typeOf reflect.Type, group *sync.WaitGroup, mutex *sync.RWMutex) {
			defer group.Done()
			jsonField := typeOf.Field(i).Tag.Get("json")
			if !strings.Contains(jsonField, "-") {
				val := values[i]
				switch v := val.(type) {
				case [16]uint8:
					p := pgtype.UUID{Bytes: v, Valid: true}
					value, err := p.MarshalJSON()
					if err != nil {
						val = err
					}
					val = strings.ReplaceAll(string(value), "\"", "")
				}
				mutex.Lock()
				if val != nil && !reflect.ValueOf(val).IsZero() {
					pointer.Elem().Field(i).Set(reflect.ValueOf(val))
				}
				mutex.Unlock()
				return
			}
		}(i, typeOf, group, mutex)
	}
	group.Wait()
	pointerInterface := pointer.Elem().Interface()
	return pointerInterface, nil
}
