package osm

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05Z"
)

// Tag represents osm tag
type Tag struct {
	K string `xml:"k,attr" json:"k"`
	V string `xml:"v,attr" json:"v"`
}

// Tags is Tag array
type Tags []*Tag

// Scan - Implement the database/sql scanner interface
func (ts *Tags) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), ts)
}

// Time is time with osm time format
type Time time.Time

func (t *Time) String() string {
	return time.Time(*t).UTC().Format(timeFormat)
}

// Scan - Implement the database/sql scanner interface
func (t *Time) Scan(v interface{}) error {
	return t.processTime(v.(string))
}

// UnmarshalJSON - implement Unmarshaler interface
func (t *Time) UnmarshalJSON(b []byte) error {
	tr := strings.Trim(string(b), "\"")
	return t.processTime(tr)
}

// MarshalXMLAttr - implement MarshalXMLAttr interface
func (t *Time) MarshalXMLAttr(name xml.Name) (attr xml.Attr, err error) {
	attr.Name = name
	attr.Value = t.String()
	return
}

// UnmarshalXMLAttr - implement UnmarshalXMLAttr interface
func (t *Time) UnmarshalXMLAttr(attr xml.Attr) error {
	return t.processTime(attr.Value)
}

func (t *Time) processTime(s string) error {
	pt, err := time.Parse(timeFormat, s)
	if err != nil {
		return err
	}
	*t = Time(pt)
	return nil
}
