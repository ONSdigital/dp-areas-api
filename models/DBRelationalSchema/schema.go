package DBRelationalSchema

// Contains relational database schema
const DBSchema = `{
    "dp-areas-api": {
        "tables": {
            "area": {
                "creation_order": 2,
                "primary_keys": "code",
                "columns": {
                    "code": {
                        "data_type": "VARCHAR(50)",
                        "constraints": ""
                    },
                    "active_from": {
                        "data_type": "TIMESTAMP",
                        "constraints": ""
                    },
                    "active_to": {
                        "data_type": "TIMESTAMP",
                        "constraints": ""
                    },
                    "geometric_area": {
                        "data_type": "VARCHAR",
                        "constraints": ""
                    },
                    "visible": {
                        "data_type": "BOOLEAN",
                        "constraints": ""
                    },
                    "area_type_id": {
                        "data_type": "INT",
                        "constraints": "REFERENCES area_type(id)"
                    }
                }
            },
            "area_name": {
                "creation_order": 4,
                "primary_keys": "id",
                "columns": {
                    "id": {
                        "data_type": "SERIAL",
                        "constraints": ""
                    },
                    "area_code": {
                        "data_type": "VARCHAR(50)",
                        "constraints": "REFERENCES area(code)"
                    },
                    "name": {
                        "data_type": "VARCHAR(50)",
                        "constraints": "UNIQUE"
                    },
                    "active_from": {
                        "data_type": "TIMESTAMP",
                        "constraints": ""
                    },
                    "active_to": {
                        "data_type": "TIMESTAMP",
                        "constraints": ""
                    }
                }
            },
            "area_relationship": {
                "creation_order": 3,
                "primary_keys": "area_code,rel_area_code",
                "columns": {
                    "area_code": {
                        "data_type": "VARCHAR(50)",
                        "constraints": "REFERENCES area(code)"
                    },
                    "rel_area_code": {
                        "data_type": "VARCHAR(50)",
                        "constraints": "REFERENCES area(code)"
                    },
                    "rel_type_id": {
                        "data_type": "INT",
                        "constraints": "REFERENCES relationship_type(id)"
                    }
                }
            },
            "area_type": {
                "creation_order": 0,
                "primary_keys": "id",
                "columns": {
                    "id": {
                        "data_type": "SERIAL",
                        "constraints": ""
                    },
                    "name": {
                        "data_type": "VARCHAR(50)",
                        "constraints": ""
                    }
                }
            },
            "relationship_type": {
                "creation_order": 1,
                "primary_keys": "id",
                "columns": {
                    "id": {
                        "data_type": "SERIAL",
                        "constraints": ""
                    },
                    "name": {
                        "data_type": "VARCHAR(50)",
                        "constraints": ""
                    }
                }
            }
        }
    }
}`
