package models

// Image represents an image metadata model as it is stored in mongoDB and json representation for API
type Image struct { //!!! re-work this to suit topic
	ID           string `bson:"_id,omitempty"           json:"id,omitempty"`
	CollectionID string `bson:"collection_id,omitempty" json:"collection_id,omitempty"`
	State        string `bson:"state,omitempty"         json:"state,omitempty"`
	Filename     string `bson:"filename,omitempty"      json:"filename,omitempty"`
	//	License      *License            `bson:"license,omitempty"       json:"license,omitempty"`
	//	Links        *ImageLinks         `bson:"links,omitempty"         json:"links,omitempty"`
	//	Upload       *Upload             `bson:"upload,omitempty"        json:"upload,omitempty"`
	Type string `bson:"type,omitempty"          json:"type,omitempty"`
	//	Downloads    map[string]Download `bson:"downloads,omitempty"     json:"-"`
}
