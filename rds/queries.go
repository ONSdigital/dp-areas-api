package rds

const (
	getBasicArea      = "select id, code, active from areas_basic where id=$1"
	getAreaCode       = "select code from area where code = $1"
	getRelationShipAreas = "select an.area_code, an.name from area_name as an, area_relationship as ar where ar.rel_area_code = an.area_code and ar.area_code =$1"
)
