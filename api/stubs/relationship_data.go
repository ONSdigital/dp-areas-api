package stubs

import "github.com/ONSdigital/dp-areas-api/models"

var Relationships = map[string][]models.AreaRelationShips{
	"E92000001": {
		{"E12000001", "North East", "/v1/area/E12000001"},
		{"E12000002", "North West", "/v1/area/E12000002"},
		{"E12000003", "Yorkshire and The Humbe", "/v1/area/E12000003"},
		{"E12000004", "East Midlands", "/v1/area/E12000004"},
		{"E12000005", "West Midlands", "/v1/area/E12000005"},
		{"E12000006", "East of England", "/v1/area/E12000006"},
		{"E12000007", "London", "/v1/area/E12000007"},
		{"E12000008", "South East", "/v1/area/E12000008"},
		{"E12000009", "South West", "/v1/area/E12000009"},
	},
	"W92000004": {
		{"W06000001", "Isle of Anglesey", "/v1/area/W06000001"},
		{"W06000002", "Gwynedd", "/v1/area/W06000002"},
		{"W06000003", "Conwy", "/v1/area/W06000003"},
		{"W06000004", "Denbighshire", "/v1/area/W06000004"},
		{"W06000005", "Flintshire", "/v1/area/W06000005"},
		{"W06000006", "Wrexham", "/v1/area/W06000006"},
		{"W06000008", "Ceredigion", "/v1/area/W06000008"},
		{"W06000009", "Pembrokeshire", "/v1/area/W06000009"},
	},
}
