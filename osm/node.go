package osm

import (
	"encoding/json"
	"encoding/xml"
)

// xmlNameJSONTypeNode is kind of a hack to encode the proper json
// object type attribute for this struct type.
type xmlNameJSONTypeNode xml.Name

func (x xmlNameJSONTypeNode) MarshalJSON() ([]byte, error) {
	return []byte(`"node"`), nil
}

// Node is an osm point and allows for marshalling to/from osm xml.
type Node struct {
	XMLName     xmlNameJSONTypeNode `xml:"node" json:"type"`
	ID          int64               `db:"id" xml:"id,attr" json:"id"`
	Lat         *float64            `db:"lat" xml:"lat,attr" json:"lat"`
	Lon         *float64            `db:"lon" xml:"lon,attr" json:"lon"`
	User        *string             `db:"user" xml:"user,attr" json:"user,omitempty"`
	UserID      *int64              `db:"uid" xml:"uid,attr" json:"uid,omitempty"`
	Visible     bool                `db:"visible" xml:"visible,attr" json:"visible"`
	Version     int                 `db:"version" xml:"version,attr" json:"version,omitempty"`
	ChangesetID int64               `db:"changeset" xml:"changeset,attr" json:"changeset,omitempty"`
	Timestamp   TimeOSM             `db:"timestamp" xml:"timestamp,attr" json:"timestamp"`
	Tags        Tags                `db:"tags" xml:"tag" json:"tags,omitempty"`
}

// ObjectID returns the object id of the node.
func (n *Node) ObjectID() int64 {
	return n.ID
}

// Nodes is a list of nodes with helper functions on top.
type Nodes []*Node

// Scan - Implement the database/sql scanner interface
func (ns *Nodes) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), ns)
}
