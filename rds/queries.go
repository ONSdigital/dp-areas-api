package rds

const (
	getBasicArea      = "select id, code, active from areas_basic where id=$1"
	getAreaCode       = "select code from area where code = $1"
	getRelationShipAreas = "select area_code, name from area_name, area_relationship where rel_area_code = area_code and area_code = $1"
)
