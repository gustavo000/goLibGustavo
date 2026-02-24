package query_builder

import (
	"fmt"
	"strings"
	"time"

	"github.com/kataras/iris/v12"

	"cust-rtmn-orch-library/constants"
	open_telemetry_span "cust-rtmn-orch-library/libs/open-telemetry-span"
	"cust-rtmn-orch-library/resources/properties"
)

const (
	QueryFrom         = " FROM %s.%s %s"
	QuerySave         = "@save_%s,"
	QueryColumnsValue = "(%s) VALUES (%s)"
	QueryOnlyValue    = "(%s)"
	QueryUpdate       = "%s = @update_%s,"
)

type Where struct {
	Field               string
	Operator            string
	Value               string
	Force               string
	Function            string
	EncapsulateFunction string
}

type Field struct {
	Column string
	As     string
	Value  any
}

type InsertResult struct {
	Id int32 `json:"id"`
}

type QueryType int

const (
	SELECT QueryType = 1
	INSERT QueryType = 2
	UPDATE QueryType = 3
	DELETE QueryType = 4
)

type query struct {
	columns    []Column
	database   string
	table      string
	schema     string
	queryParts []string
	modelCount int
	queryArgs  []any
	tableChar  string
	tableCount int
	ctx        iris.Context
	queryType  QueryType
	namedArgs  pgx.NamedArgs
	model      any
	firstOnly  bool
}

func QueryBuilder(ctx iris.Context) *query {
	q := &query{ctx: ctx}
	q.tableCount = 1
	q.tableChar = strings.ToLower(string(toChar(q.tableCount)))
	return q
}

func (q *query) WithDatabase(database string) *query {
	switch properties.GetProperty().GetEnv() {
	case "TEST":
		database += "-uat"
	case "BETA":
		database += "-beta"
	case "PROD":
		database += "-prod"
	}
	q.database = database
	return q
}

func (q *query) WithSchema(schema string) *query {
	q.schema = schema
	return q
}

func (q *query) WithTable(table string) *query {
	q.table = table
	return q
}

func (q *query) WithModel(model any) *query {
	q.model = model
	return q
}

func (q *query) AddNamedArgs(namedArgs pgx.NamedArgs) *query {
	for key, value := range namedArgs {
		q.setNamedArg("custom", key, value)
	}
	return q
}

type QueryOption func(*query)

func WithStruct(object any) QueryOption {
	return func(q *query) {
		q.modelCount++
		if q.model == nil {
			q.model = object
		}
		var qrs string
		columns := q.generateColumns(object)
		switch q.queryType {
		case SELECT:
			for _, column := range columns {
				column.JsonField = strings.ReplaceAll(column.JsonField, ",omitempty", "")
				qrs += fmt.Sprintf("%s.%s as %s,", q.tableChar, column.Name, column.JsonField)
			}
			qrs = strings.TrimSuffix(qrs, ",")
			qrs += fmt.Sprintf(QueryFrom, q.schema, q.table, q.tableChar)
			break
		case INSERT:
			var columnsQuery string
			var valuesQuery string
			for key, value := range q.generateMap(columns) {
				columnsQuery += fmt.Sprintf("%s,", key)
				valuesQuery += fmt.Sprintf(QuerySave, key)
				q.setNamedArg("save", key, value)
			}
			valuesQuery = strings.TrimSuffix(valuesQuery, ",")
			columnsQuery = strings.TrimSuffix(columnsQuery, ",")
			if q.modelCount == 1 {
				qrs = fmt.Sprintf(QueryColumnsValue, columnsQuery, valuesQuery)
			} else {
				qrs = fmt.Sprintf(QueryOnlyValue, valuesQuery)
			}
			break
		case UPDATE:
			qrs += "SET "
			for key, value := range q.generateMap(columns) {
				qrs += fmt.Sprintf(QueryUpdate, key, key)
				q.setNamedArg("update", key, value)
			}
			qrs = strings.TrimSuffix(qrs, ",")
			break
		}
		q.queryParts = append(q.queryParts, qrs)
	}
}

func WithMap(object map[string]any) QueryOption {
	return func(q *query) {
		var qrs string
		switch q.queryType {
		case SELECT:
			for key, value := range object {
				qrs += fmt.Sprintf("%s.%s as %s,", q.tableChar, key, value)
			}
			qrs = strings.TrimSuffix(qrs, ",")
			qrs += fmt.Sprintf(QueryFrom, q.schema, q.table, q.tableChar)
			break
		case INSERT:
			var columnsQuery string
			var valuesQuery string
			for key, value := range object {
				columnsQuery += fmt.Sprintf("%s,", key)
				valuesQuery += fmt.Sprintf(QuerySave, key)
				q.setNamedArg("save", key, value)
			}
			valuesQuery = strings.TrimSuffix(valuesQuery, ",")
			columnsQuery = strings.TrimSuffix(columnsQuery, ",")
			qrs += fmt.Sprintf(QueryColumnsValue, columnsQuery, valuesQuery)
			break
		case UPDATE:
			qrs += "SET "
			for key, value := range object {
				qrs += fmt.Sprintf(QueryUpdate, key, key)
				q.setNamedArg("update", key, value)
			}
			qrs = strings.TrimSuffix(qrs, ",")
			break
		}
		q.queryParts = append(q.queryParts, qrs)
	}
}

func WithFields(fields ...Field) QueryOption {
	return func(q *query) {
		var qrs string
		switch q.queryType {
		case SELECT:
			for _, field := range fields {
				qrs += fmt.Sprintf("%s.%s", q.tableChar, field.Column)
				if field.As != "" {
					qrs += fmt.Sprintf("as %s", field.As)
				}
				qrs += ","
			}
			qrs = strings.TrimSuffix(qrs, ",")
			qrs += fmt.Sprintf(QueryFrom, q.schema, q.table, q.tableChar)
			break
		case INSERT:
			var columnsQuery string
			var valuesQuery string
			for _, field := range fields {
				columnsQuery += fmt.Sprintf("%s,", field.Column)
				valuesQuery += fmt.Sprintf(QuerySave, field.Column)
				q.setNamedArg("save", field.Column, field.Value)
			}
			valuesQuery = strings.TrimSuffix(valuesQuery, ",")
			columnsQuery = strings.TrimSuffix(columnsQuery, ",")
			qrs += fmt.Sprintf(QueryColumnsValue, columnsQuery, valuesQuery)
			break
		case UPDATE:
			qrs += "SET "
			for _, field := range fields {
				qrs += fmt.Sprintf(QueryUpdate, field.Column, field.Column)
				q.setNamedArg("update", field.Column, field.Value)
			}
			qrs = strings.TrimSuffix(qrs, ",")
			break
		}
		q.queryParts = append(q.queryParts, qrs)
	}
}

func (q *query) SelectFirst(args ...QueryOption) *query {
	q.queryType = SELECT
	q.queryParts = []string{
		"SELECT",
	}
	for _, arg := range args {
		arg(q)
	}
	q.firstOnly = true
	return q
}

func (q *query) Select(args ...QueryOption) *query {
	q.queryType = SELECT
	q.queryParts = []string{
		"SELECT",
	}
	for _, arg := range args {
		arg(q)
	}
	return q
}

func (q *query) Insert(args ...QueryOption) *query {
	q.queryType = INSERT
	q.queryParts = []string{
		"INSERT INTO",
	}
	q.queryParts = append(q.queryParts, fmt.Sprintf("%s.%s", q.schema, q.table))
	for _, arg := range args {
		arg(q)
	}
	return q
}

func (q *query) Update(args ...QueryOption) *query {
	q.queryType = UPDATE
	q.queryParts = []string{
		"UPDATE",
	}
	q.queryParts = append(q.queryParts, fmt.Sprintf("%s.%s", q.schema, q.table))
	for _, arg := range args {
		arg(q)
	}
	return q
}

func (q *query) Delete(args ...QueryOption) *query {
	q.queryType = DELETE
	q.queryParts = []string{
		"DELETE",
	}
	q.queryParts = append(q.queryParts, fmt.Sprintf("FROM %s.%s", q.schema, q.table))
	for _, arg := range args {
		arg(q)
	}
	return q
}

type WhereOption func(*query)

func (q *query) setNamedArg(typeArg string, field string, value any) {
	if q.namedArgs == nil {
		q.namedArgs = make(pgx.NamedArgs)
	}
	switch s := value.(type) {
	case string:
		q.namedArgs[typeArg+"_"+field] = strings.ReplaceAll(s, "'", " ")
	default:
		q.namedArgs[typeArg+"_"+field] = s
	}
}

func IsLike(field string, value any) WhereOption {
	return func(q *query) {
		q.queryParts = append(q.queryParts, fmt.Sprintf("%s LIKE @where_%d_%s", field, len(q.queryArgs), field))
		q.setNamedArg("where", fmt.Sprintf("%d_%s", len(q.queryArgs), field), value)
		q.queryArgs = append(q.queryArgs, value)
	}
}

func IsEqual(field string, value any) WhereOption {
	return func(q *query) {
		q.queryParts = append(q.queryParts, fmt.Sprintf("%s = @where_%d_%s", field, len(q.queryArgs), field))
		q.setNamedArg("where", fmt.Sprintf("%d_%s", len(q.queryArgs), field), value)
		q.queryArgs = append(q.queryArgs, value)
	}
}

func ContainsInArray(field string, values ...string) WhereOption {
	return func(q *query) {
		q.queryParts = append(q.queryParts, fmt.Sprintf("%s @> @where_%d_%s", field, len(q.queryArgs), field))
		q.setNamedArg("where", fmt.Sprintf("%d_%s", len(q.queryArgs), field), values)
		result := "'" + strings.Join(values, "','") + "'"
		q.queryArgs = append(q.queryArgs, fmt.Sprintf("Array[%s]", result))
	}
}

func NotEqual(field string, value any) WhereOption {
	return func(q *query) {
		q.queryParts = append(q.queryParts, fmt.Sprintf("%s != @where_%d_%s", field, len(q.queryArgs), field))
		q.setNamedArg("where", fmt.Sprintf("%d_%s", len(q.queryArgs), field), value)
		q.queryArgs = append(q.queryArgs, value)
	}
}

func IsGreaterThan(field string, value any) WhereOption {
	return func(q *query) {
		q.queryParts = append(q.queryParts, fmt.Sprintf("%s > @where_%d_%s", field, len(q.queryArgs), field))
		q.setNamedArg("where", fmt.Sprintf("%d_%s", len(q.queryArgs), field), value)
		q.queryArgs = append(q.queryArgs, value)
	}
}

func IsLessThan(field string, value any) WhereOption {
	return func(q *query) {
		q.queryParts = append(q.queryParts, fmt.Sprintf("%s < @where_%d_%s", field, len(q.queryArgs), field))
		q.setNamedArg("where", fmt.Sprintf("%d_%s", len(q.queryArgs), field), value)
		q.queryArgs = append(q.queryArgs, value)
	}
}

func IsGreaterOrEqualThan(field string, value any) WhereOption {
	return func(q *query) {
		q.queryParts = append(q.queryParts, fmt.Sprintf("%s >= @where_%d_%s", field, len(q.queryArgs), field))
		q.setNamedArg("where", fmt.Sprintf("%d_%s", len(q.queryArgs), field), value)
		q.queryArgs = append(q.queryArgs, value)
	}
}

func IsLessOrEqualThan(field string, value any) WhereOption {
	return func(q *query) {
		q.queryParts = append(q.queryParts, fmt.Sprintf("%s <= @where_%d_%s", field, len(q.queryArgs), field))
		q.setNamedArg("where", fmt.Sprintf("%d_%s", len(q.queryArgs), field), value)
		q.queryArgs = append(q.queryArgs, value)
	}
}

func And(wheres ...WhereOption) WhereOption {
	return func(q *query) {
		q.queryParts = append(q.queryParts, "AND")
		for _, where := range wheres {
			where(q)
		}
	}
}

func Or(wheres ...WhereOption) WhereOption {
	return func(q *query) {
		q.queryParts = append(q.queryParts, "OR")
		for _, where := range wheres {
			where(q)
		}
	}
}

func (q *query) Where(wheres ...WhereOption) *query {
	q.queryParts = append(q.queryParts, "WHERE")
	for _, where := range wheres {
		where(q)
	}
	return q
}

func (q *query) Returning(field string) *query {
	q.queryParts = append(q.queryParts, "RETURNING "+field)
	return q
}

func (q *query) Perform() (any, error) {
	spanChild := open_telemetry_span.StartChildSpan(q.ctx)
	defer open_telemetry_span.EndChildSpan(spanChild)
	resultQuery := strings.Join(q.queryParts, " ") + ";"
	q.setOnContext(constants.QUERY_SYNTAX, resultQuery)
	queryJobTime := time.Now()
	rows, conn, err := q.executeQuery()
	if conn != nil {
		defer conn.Release()
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var resultList []any
	var rowCount = 0
	for rows.Next() {
		res, errRead := q.getResultFromRow(rows)
		if errRead != nil {
			return nil, errRead
		}
		resultList = append(resultList, res)
		rowCount++
	}

	q.setOnContext(constants.QUERY_JOB, time.Since(queryJobTime))
	if len(resultList) == 0 {
		return nil, fmt.Errorf("query without results")
	} else if q.firstOnly {
		return resultList[0], nil
	}
	return resultList, nil
}
