package rds

import "fmt"

const (
	getBasicArea                      = "select id, code, active from areas_basic where id=$1"
	getAreaCode                       = "select code from area where code = $1"
	getAreaType                       = "select id from area_type where name = $1"
	getRelationShipAreas              = "select an.area_code, an.name from area_name as an, area_relationship as ar where ar.rel_area_code = an.area_code and ar.area_code =$1"
	upsertAreaName                    = "insert into area_name(area_code, name, active_from, active_to) values($1, $2, $3, $4) on conflict(name) do update set active_from=$3,active_to=$4"
	insertArea                        = "insert into area(code, active_from, active_to, geometric_area, area_type_id) values($1, $2, $3, $4, $5)"
	updateAreaOnConflict              = "on conflict(code) do update set active_from=$2, active_to=$3,geometric_area=$4,area_type_id=$5 returning (xmax = 0) as inserted"
	areaTypeInsertTransaction         = "insert into area_type(name) select $1 where not exists (select * from area_type where name = $2)"
	areaInsertTransaction             = "insert into area(code, active_from, active_to, area_type_id, geometric_area) VALUES($1, $2, $3, $4, $5) on conflict (code) do nothing"
	relationshipTypeInsertTransaction = "insert into relationship_type(name) select $1 where not exists (select * from relationship_type where name = $2)"

	areaNameInsertTransaction         = "insert into area_name(area_code, name, active_from, active_to) VALUES($1, $2, $3, $4)"
	areaRelationshipInsertTransaction = "insert into area_relationship(area_code, rel_area_code, rel_type_id) VALUES($1, $2, $3)"
)

var upsertArea = fmt.Sprintf("%s %s", insertArea, updateAreaOnConflict)
