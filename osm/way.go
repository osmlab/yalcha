package osm

import (
	"encoding/json"
	"encoding/xml"
)

// xmlNameJSONTypeWay is kind of a hack to encode the proper json
// object type attribute for this struct type.
type xmlNameJSONTypeWay xml.Name

func (x xmlNameJSONTypeWay) MarshalJSON() ([]byte, error) {
	return []byte(`"way"`), nil
}

// Way is an osm way and allows for marshalling to/from osm xml.
type Way struct {
	XMLName     xmlNameJSONTypeWay `xml:"way" json:"type"`
	ID          int64              `db:"id" xml:"id,attr" json:"id"`
	Visible     bool               `db:"visible" xml:"visible,attr" json:"visible"`
	Version     int                `db:"version" xml:"version,attr" json:"version,omitempty"`
	User        *string            `db:"user" xml:"user,attr,omitempty" json:"user,omitempty"`
	UserID      *int64             `db:"uid" xml:"uid,attr,omitempty" json:"uid,omitempty"`
	ChangesetID int64              `db:"changeset" xml:"changeset,attr" json:"changeset,omitempty"`
	Timestamp   Time               `db:"timestamp" xml:"timestamp,attr" json:"timestamp"`
	Nodes       wayNodes           `db:"nodes" xml:"nd" json:"nodes"`
	Tags        Tags               `db:"tags" xml:"tag" json:"tags,omitempty"`
}

// wayNodes represents a collection of way nodes.
type wayNodes []wayNode

// Scan - Implement the database/sql scanner interface
func (wn *wayNodes) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &wn)
}

// wayNode is a short node used as part of ways and relations in the osm xml.
type wayNode struct {
	ID int64 `xml:"ref,attr"`
}

func (wn *wayNode) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &wn.ID)
}

// ObjectID returns the object id of the way.
func (w *Way) ObjectID() int64 {
	return w.ID
}

// Ways is a list of ways with helper functions on top.
type Ways []*Way

// Scan - Implement the database/sql scanner interface
func (ws *Ways) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), ws)
}
