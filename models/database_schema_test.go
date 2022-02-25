package models_test

import (
	"testing"

	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-areas-api/models/DBRelationalSchema"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	area_query              = "CREATE TABLE IF NOT EXISTS area (PRIMARY KEY (code), active_from TIMESTAMP , active_to TIMESTAMP , area_type_id INT REFERENCES area_type(id), code VARCHAR(50) , geometric_area VARCHAR , visible BOOLEAN )"
	area_type_query         = "CREATE TABLE IF NOT EXISTS area_type (PRIMARY KEY (id), id SERIAL , name VARCHAR(50) )"
	area_name_query         = "CREATE TABLE IF NOT EXISTS area_name (PRIMARY KEY (id), active_from TIMESTAMP , active_to TIMESTAMP , area_code VARCHAR(50) REFERENCES area(code), id SERIAL , name VARCHAR(50) UNIQUE)"
	relationship_type_query = "CREATE TABLE IF NOT EXISTS relationship_type (PRIMARY KEY (id), id SERIAL , name VARCHAR(50) )"
	area_relationship_query = "CREATE TABLE IF NOT EXISTS area_relationship (PRIMARY KEY (area_code,rel_area_code), area_code VARCHAR(50) REFERENCES area(code), rel_area_code VARCHAR(50) REFERENCES area(code), rel_type_id INT REFERENCES relationship_type(id))"
)

func TestSetup(t *testing.T) {
	Convey("Ensure database schema model is built correctly", t, func() {
		Convey("When a valid schema string is used - schema model built successfully", func() {
			databaseSchema := models.DatabaseSchema{
				DBName:       "dp-areas-api",
				SchemaString: DBRelationalSchema.DBSchema,
			}
			err := databaseSchema.BuildDatabaseSchemaModel()
			So(err, ShouldEqual, nil)
			So(len(databaseSchema.Tables), ShouldBeGreaterThan, 1)
			// sample from built schema model
			So(databaseSchema.Tables["area"]["creation_order"].(float64), ShouldEqual, 2)
			So(databaseSchema.Tables["area"]["primary_keys"].(string), ShouldEqual, "code")
			So(len(databaseSchema.Tables["area"]["columns"].(map[string]interface{})), ShouldEqual, 6)
		})

		Convey("When an invalid schema string is used - error generated", func() {
			databaseSchema := models.DatabaseSchema{
				DBName: "dp-areas-api",
				SchemaString: `{
					"dp-areas-api": {
						"tables": {
							"area": {
								"creation_order": 2,
								"primary_keys": "code",
								"columns": {,
				}`,
			}
			err := databaseSchema.BuildDatabaseSchemaModel()
			So(err, ShouldNotEqual, nil)
			So(len(databaseSchema.Tables), ShouldEqual, 0)
			So(err.Error(), ShouldEqual, "invalid character ',' looking for beginning of object key string")
		})

		Convey("When an no schema string supplied - error generated", func() {
			databaseSchema := models.DatabaseSchema{
				DBName:       "dp-areas-api",
				SchemaString: `""`,
			}
			err := databaseSchema.BuildDatabaseSchemaModel()
			So(err, ShouldNotEqual, nil)
			So(len(databaseSchema.Tables), ShouldEqual, 0)
			So(err.Error(), ShouldEqual, "json: cannot unmarshal string into Go value of type map[string]models.DatabaseSchema")
		})

		Convey("Ensure execution list built successfully", func() {
			databaseSchema := models.DatabaseSchema{
				DBName:       "dp-areas-api",
				SchemaString: DBRelationalSchema.DBSchema,
			}
			_ = databaseSchema.BuildDatabaseSchemaModel()
			databaseSchema.TableSchemaBuilder()
			So(len(databaseSchema.ExecutionList), ShouldBeGreaterThanOrEqualTo, 5)
			So(databaseSchema.ExecutionList[0], ShouldEqual, area_type_query)
			So(databaseSchema.ExecutionList[1], ShouldEqual, relationship_type_query)
			So(databaseSchema.ExecutionList[2], ShouldEqual, area_query)
			So(databaseSchema.ExecutionList[3], ShouldEqual, area_relationship_query)
			So(databaseSchema.ExecutionList[4], ShouldEqual, area_name_query)
		})
	})
}
