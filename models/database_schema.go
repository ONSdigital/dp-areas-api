package models

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var (
	CreateTableQuery = "CREATE TABLE IF NOT EXISTS %s (\n\tPRIMARY KEY (%s),\n%s)"
)

// DatabaseSchema database schema model
type DatabaseSchema struct {
	DBName,
	SchemaString  string
	ExecutionList []string
	Tables        map[string]map[string]interface{}
}

// BuildDatabaseSchemaModel build db schema model
func (db *DatabaseSchema) BuildDatabaseSchemaModel() error {
	dbSchemaData := make(map[string]DatabaseSchema, 1)
	str, err := db.CleanSchemaString(db.SchemaString)
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
	//as we're apply a FK constraint, order of creation is important, so build array of fixed size to preserve
	db.ExecutionList = make([]string, len(db.Tables))
	for table := range db.Tables {
		var (
			schemaHandleCols = db.Tables[table]["columns"].(map[string]interface{})
			columnData = ""
			terminator = ","
			columnCount = 1
		)
		for name, data := range schemaHandleCols {
			d := data.(map[string]interface{})
			if columnCount == len(schemaHandleCols){
				terminator = ""
			}
			columnData = columnData+fmt.Sprintf("\t%s %s %s%s\n", name, d["data_type"].(string), d["constraints"].(string), terminator)
			db.ExecutionList[int(db.Tables[table]["creation_order"].(float64))] = fmt.Sprintf(CreateTableQuery, table, db.Tables[table]["primary_keys"], columnData)
			columnCount++
		}
	}
}

// cleanSchemaString cleans db schema string
func (db *DatabaseSchema) CleanSchemaString(schemaString string) (*string, error) {
	reg, err := regexp.Compile("[\n\t]")
	if err != nil {
		return nil, err
	}
	cleanSchemaString := reg.ReplaceAllString(schemaString, "")
	return &cleanSchemaString, nil
}
