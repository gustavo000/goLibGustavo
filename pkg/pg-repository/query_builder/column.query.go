package query_builder

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/gustavo000/goLibGustavo/pkg/functions"
)

type Column struct {
	Name        string              `json:"Column"`
	Value       any                 `json:"value"`
	Type        string              `json:"type"`
	StructField reflect.StructField `json:"structField"`
	TypeOf      reflect.Type        `json:"typeOf"`
	JsonField   string              `json:"jsonField"`
	Generate    bool                `json:"generate"`
	Require     bool                `json:"require"`
	PrimaryKey  bool                `json:"pk"`
	ParseTo     string              `json:"parseTo"`
	Serial      bool                `json:"serial"`
}

func GetColumn(field reflect.StructField, object map[string]any) Column {
	var column = Column{}
	column.StructField = field
	column.JsonField = field.Tag.Get("json")
	pg := field.Tag.Get("pg")
	if pg != "" {
		mapStrings := make(map[string]any)
		options := strings.Split(pg, ",")
		for _, option := range options {
			sOption := strings.Split(option, "=")
			switch sOption[0] {
			case "require", "pk", "generate", "serial":
				if len(sOption) > 1 {
					mapStrings[sOption[0]] = sOption[1] == "true"
				} else {
					mapStrings[sOption[0]] = true
				}
			default:
				mapStrings[sOption[0]] = sOption[1]
			}

		}
		err := functions.ParseTo(&mapStrings, &column)
		if err != nil {
			return column
		}
		if value, ok := object[column.JsonField]; ok {
			switch column.ParseTo {
			case "timestamp":
				switch s := value.(type) {
				case string:
					parse, err := time.Parse(time.RFC3339, strings.ReplaceAll(s, "\"", ""))
					if err == nil {
						column.Value = parse.Format("2006-01-02 15:04:05")
					}
				default:
					buf, _ := json.Marshal(s)
					parse, err := time.Parse(time.RFC3339, strings.ReplaceAll(string(buf), "\"", ""))
					if err == nil {
						column.Value = parse.Format("2006-01-02 15:04:05")
					}
				}
			case "string":
				marshal, err := json.Marshal(value)
				if err == nil {
					column.Value = string(marshal)
				}
			default:
				column.Value = value
			}
			column.TypeOf = reflect.TypeOf(value)
		}
	}
	return column
}
