package models

// Topic represents topic schema as it is stored in mongoDB
// and is used for marshaling and unmarshaling json representation for API
type Topic struct {
	ID          string   `bson:"_id,omitempty"          json:"id,omitempty"`
	Description string   `bson:"description,omitempty"  json:"description,omitempty"`
	Title       string   `bson:"title,omitempty"        json:"title,omitempty"`
	Keywords    []string `bson:"keywords,omitempty"     json:"keywords,omitempty"`
	State       string   `bson:"state,omitempty"        json:"state,omitempty"`
}
