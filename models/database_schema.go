package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

var (
	CreateTableQuery = "CREATE TABLE IF NOT EXISTS %s (PRIMARY KEY (%s), %s)"
)

// DatabaseSchema database schema model
type DatabaseSchema struct {
	DBName,
	SchemaString string
	ExecutionList []string
	Tables        map[string]map[string]interface{}
}

// BuildDatabaseSchemaModel build db schema model
func (db *DatabaseSchema) BuildDatabaseSchemaModel() error {
	dbSchemaData := make(map[string]DatabaseSchema, 1)
	str, err := cleanString(db.SchemaString)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(*str), &dbSchemaData)
	if err != nil {
		return err
	}
	db.Tables = dbSchemaData[db.DBName].Tables
	return nil
}

// TableSchemaBuilder builds table schema
func (db *DatabaseSchema) TableSchemaBuilder() {
	// as we're applied a FK constraint, order of creation is important, so build array of fixed size to preserve
	db.ExecutionList = make([]string, len(db.Tables))
	for table := range db.Tables {
		var (
			schemaHandleCols = db.Tables[table]["columns"].(map[string]interface{})
			colArray         = make([]string, 0)
		)
		for name, data := range schemaHandleCols {
			d := data.(map[string]interface{})
			colString := fmt.Sprintf("%s %s %s", name, d["data_type"].(string), d["constraints"].(string))
			colArray = append(colArray, colString)
		}
		sort.Strings(colArray)
		columnData := strings.Join(colArray, ", ")
		db.ExecutionList[int(db.Tables[table]["creation_order"].(float64))] = fmt.Sprintf(CreateTableQuery, table, db.Tables[table]["primary_keys"], columnData)
	}
}

// cleanString cleans strings of unwanted chars matched in localised char set
func cleanString(str string) (*string, error) {
	reg := regexp.MustCompile("[\n\t]")
	if reg == nil {
		return nil, errors.New("failed to find any matching pattern")
	}
	cleanSchemaString := reg.ReplaceAllString(str, "")
	return &cleanSchemaString, nil
}
