package rds

const (
	getBasicArea                      = "select id, code, active from areas_basic where id=$1"
	getAreaCode                       = "select code from area where code = $1"
	getRelationShipAreas              = "select an.area_code, an.name from area_name as an, area_relationship as ar where ar.rel_area_code = an.area_code and ar.area_code =$1"
	getRelationShipAreasWithParameter = "select an.area_code, an.name from area_name as an, area_relationship as ar where ar.rel_area_code = an.area_code and ar.area_code = $1 and ar.rel_type_id = (select id from relationship_type where name = $2)"
	getAncestors                      = "select an.area_code, an.name from area_name as an, area_relationship as ar where ar.area_code = an.area_code and ar.rel_area_code = $1 and ar.rel_type_id = (select id from relationship_type where name = $2)"
	areaTypeInsertTransaction         = "insert into area_type(name) select $1 where not exists (select * from area_type where name = $2)"
	areaInsertTransaction             = "insert into area(code, active_from, active_to, area_type_id, geometric_area) VALUES($1, $2, $3, $4, $5) on conflict (code) do nothing"
)
