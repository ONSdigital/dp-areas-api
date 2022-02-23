package DBRelationalData

// var representing test relationship_type data
var RelationshipTypeData = map[string]map[string]interface{}{
	"child": {
		"creation_order": 0,
		"columns":        "name",
		"values":         "child",
	},
	"bordering": {
		"creation_order": 1,
		"columns":        "name",
		"values":         "bordering",
	},
	"supercedes": {
		"creation_order": 2,
		"columns":        "name",
		"values":         "supercedes",
	},
	"superceded_by": {
		"creation_order": 3,
		"columns":        "name",
		"values":         "superceded_by",
	},
	"statistical_neighbour": {
		"creation_order": 4,
		"columns":        "name",
		"values":         "statistical_neighbour",
	},
	"related": {
		"creation_order": 4,
		"columns":        "name",
		"values":         "related",
	},
}
