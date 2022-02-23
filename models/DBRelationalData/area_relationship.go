package DBRelationalData

// var representing test area relationship data
var AreaRelationshipData = map[string]map[string]interface{}{
	"E92000001": {
		"columns": "area_code, rel_area_code, rel_type_id",
		"values": map[string]interface{}{
			"area_code":     "E92000001",
			"rel_area_code": "E12000003",
			"rel_type_id":   "1",
		},
	},
	"W92000004": {
		"columns": "area_code, rel_area_code, rel_type_id",
		"values": map[string]interface{}{
			"area_code":     "W92000004",
			"rel_area_code": "W37000382",
			"rel_type_id":   "1",
		},
	},
	"E12000003": {
		"columns": "area_code, rel_area_code, rel_type_id",
		"values": map[string]interface{}{
			"area_code":     "E12000003",
			"rel_area_code": "E08000019",
			"rel_type_id":   "1",
		},
	},
	"W37000382": {
		"columns": "area_code, rel_area_code, rel_type_id",
		"values": map[string]interface{}{
			"area_code":     "W37000382",
			"rel_area_code": "W38000028",
			"rel_type_id":   "1",
		},
	},
}
