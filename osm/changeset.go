package osm

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// xmlNameJSONTypeCS is kind of a hack to encode the proper json
// object type attribute for this struct type.
type xmlNameJSONTypeCS xml.Name

func (x xmlNameJSONTypeCS) MarshalJSON() ([]byte, error) {
	return []byte(`"changeset"`), nil
}

// A Changeset is a set of metadata around a set of osm changes.
type Changeset struct {
	XMLName       xmlNameJSONTypeCS    `xml:"changeset" json:"type"`
	ID            int64                `xml:"id,attr" json:"id"`
	User          *string              `xml:"user,attr" json:"user,omitempty"`
	UserID        *int64               `xml:"uid,attr" json:"uid,omitempty"`
	CreatedAt     Time                 `xml:"created_at,attr" json:"created_at"`
	ClosedAt      Time                 `xml:"closed_at,attr" json:"closed_at"`
	Open          bool                 `xml:"open,attr" json:"open"`
	ChangesCount  int                  `xml:"num_changes,attr,omitempty" json:"num_changes,omitempty"`
	MinLat        float64              `xml:"min_lat,attr" json:"min_lat,omitempty"`
	MaxLat        float64              `xml:"max_lat,attr" json:"max_lat,omitempty"`
	MinLon        float64              `xml:"min_lon,attr" json:"min_lon,omitempty"`
	MaxLon        float64              `xml:"max_lon,attr" json:"max_lon,omitempty"`
	CommentsCount int                  `xml:"comments_count,attr" json:"comments_count,omitempty"`
	Tags          Tags                 `xml:"tag" json:"tags,omitempty"`
	Discussion    *ChangesetDiscussion `xml:"discussion,omitempty" json:"discussion,omitempty"`
}

// ObjectID returns the object id of the changeset.
func (c *Changeset) ObjectID() int64 {
	return c.ID
}

// Changesets is a collection with some helper functions attached.
type Changesets []*Changeset

// ChangesetDiscussion is a conversation about a changeset.
type ChangesetDiscussion struct {
	Comments []*ChangesetComment `xml:"comment" json:"comments"`
}

// Scan - Implement the database/sql scanner interface
func (cd *ChangesetDiscussion) Scan(value interface{}) error {
	fmt.Println(value)
	return json.Unmarshal(value.([]byte), cd.Comments)
}

// ChangesetComment is a specific comment in a changeset discussion.
type ChangesetComment struct {
	User      string `xml:"user,attr" json:"f2"`
	UserID    int64  `xml:"uid,attr" json:"f1"`
	Timestamp Time   `xml:"date,attr" json:"f4"`
	Text      string `xml:"text" json:"f3"`
}

// Scan - Implement the database/sql scanner interface
func (cc *ChangesetComment) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), cc)
}
