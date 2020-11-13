package models

// Topic represents topic !!! fix this !!! model as it is stored in mongoDB and json representation for API
type Topic struct {
	ID          string   `bson:"_id,omitempty"          json:"id,omitempty"`
	Description string   `bson:"description,omitempty"  json:"description,omitempty"`
	Title       string   `bson:"title,omitempty"        json:"title,omitempty"`
	Keywords    []string `bson:"keywords,omitempty"     json:"keywords,omitempty"`
	State       string   `bson:"state,omitempty"        json:"state,omitempty"`
}
