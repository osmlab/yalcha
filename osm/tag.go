package osm

import "encoding/json"

// Tag represents osm tag
type Tag struct {
	K string `xml:"k,attr" json:"k"`
	V string `xml:"v,attr" json:"v"`
}

// Tags is Tag array
type Tags []Tag

// Scan - Implement the database/sql scanner interface
func (ts *Tags) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), ts)
}
