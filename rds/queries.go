package rds

const (
	getBasicArea         = "select id, code, active from areas_basic where id=$1"
	getAreaCode          = "select code from area where code = $1"
	getAreaType          = "select id from area_type where name = $1"
	getRelationShipAreas = "select an.area_code, an.name from area_name as an, area_relationship as ar where ar.rel_area_code = an.area_code and ar.area_code =$1"
	insertArea           = "insert into area(code, active_from, active_to, geometric_area, area_type_id) values($1, $2, $3, $4, $5)"
	updateArea           = "update area set active_from=$2, active_to=$3,geometric_area=$4,area_type_id=$5 where code=$1"
	upsertAreaName       = "insert into area_name(area_code, name, active_from, active_to) values($1, $2, $3, $4) on conflict set name=$2, active_from=$3,active_to=$4"
	areaTypeInsertTransaction = "insert into area_type(name) select $1 where not exists (select * from area_type where name = $2)"
	areaInsertTransaction = "insert into area(code, active_from, active_to, area_type_id, geometric_area) VALUES($1, $2, $3, $4, $5) on conflict (code) do nothing"
)
