package rds

import "fmt"

const (
	getArea = `select code, area_name.name, geometric_area, visible, area_type.name
               from area
               left join area_name on area.code = area_name.area_code
               left join area_type on area.area_type_id = area_type.id
               where code = $1`
    getAreaCode                       = "select code from area where code = $1"
    getAreaType                       = "select id from area_type where name = $1"
    getRelationShipAreas              = "select an.area_code, an.name from area_name as an, area_relationship as ar where ar.rel_area_code = an.area_code and ar.area_code =$1"
    getRelationShipAreasWithParameter = "select an.area_code, an.name from area_name as an, area_relationship as ar where ar.rel_area_code = an.area_code and ar.area_code = $1 and ar.rel_type_id = (select id from relationship_type where name = $2)"
    upsertAreaName                    = "insert into area_name(area_code, name, active_from, active_to) values($1, $2, $3, $4) on conflict(name) do update set active_from=$3,active_to=$4"
    insertArea                        = "insert into area(code, active_from, active_to, geometric_area, area_type_id, visible, land_hectares) values($1, $2, $3, $4, $5, $6, $7)"
    updateAreaOnConflict              = "on conflict(code) do update set active_from=$2, active_to=$3,geometric_area=$4,area_type_id=$5, visible=$6, land_hectares=$7 returning (xmax = 0) as inserted"
    areaTypeInsertTransaction         = "insert into area_type(name) select $1 where not exists (select * from area_type where name = $2)"
    areaInsertTransaction             = `insert into area(code, active_from, active_to, area_type_id, geometric_area, visible)
                                 VALUES($1, $2, $3, $4, $5, $6)
                                 on conflict (code) do update
                                 set active_from=$2,active_to=$3, area_type_id=$4,geometric_area=$5`
    relationshipTypeInsertTransaction = "insert into relationship_type(name) select $1 where not exists (select * from relationship_type where name = $2)"

    areaNameInsertTransaction         = "insert into area_name(area_code, name, active_from, active_to) VALUES($1, $2, $3, $4)"
	areaRelationshipInsertTransaction = "insert into area_relationship(area_code, rel_area_code, rel_type_id) VALUES($1, $2, $3) on conflict(area_code, rel_area_code) do update set rel_type_id = $3"
	getRelationShipId                 = "select id from relationship_type where name = 'child'"
    getAncestors                      = "with recursive ancestors as (select area_code from area_relationship where rel_area_code = $1 and rel_type_id = (select id from relationship_type where name = 'child') union select ar.area_code from area_relationship ar inner join ancestors a on a.area_code = ar.rel_area_code ) select a.area_code, an.name from ancestors as a, area_name as an where a.area_code = an.area_code"
)

var upsertArea = fmt.Sprintf("%s %s", insertArea, updateAreaOnConflict)
