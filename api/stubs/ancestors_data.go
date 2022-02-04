package stubs

import "github.com/ONSdigital/dp-areas-api/models"

// Ancestors is stub data for the AreasAncestors model.
var Ancestors = map[string][]models.AreasAncestors{
	"E92000001": {},
	"W92000004": {},
	"E34002743": {
		{"E92000001", "England"},
	},
}
