package osm

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"time"
)

const osmTimeFormat = "2006-01-02T15:04:05"

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

// TimeOSM is time with osm time format
type TimeOSM time.Time

func (t *TimeOSM) String() string {
	return time.Time(*t).Format(osmTimeFormat)
}

// Scan - Implement the database/sql scanner interface
func (t *TimeOSM) Scan(v interface{}) error {
	*t = TimeOSM(v.(time.Time))
	return nil
}

// UnmarshalJSON - implement Unmarshaler interface
func (t *TimeOSM) UnmarshalJSON(b []byte) error {
	tr := strings.Trim(string(b), "\"")
	return t.processTimeOSM(tr)
}

// MarshalXMLAttr - implement MarshalXMLAttr interface
func (t *TimeOSM) MarshalXMLAttr(name xml.Name) (attr xml.Attr, err error) {
	attr.Name = name
	attr.Value = t.String()
	return
}

// UnmarshalXMLAttr - implement UnmarshalXMLAttr interface
func (t *TimeOSM) UnmarshalXMLAttr(attr xml.Attr) error {
	return t.processTimeOSM(attr.Value)
}

func (t *TimeOSM) processTimeOSM(s string) error {
	pt, err := time.Parse(osmTimeFormat, s)
	if err != nil {
		return err
	}
	*t = TimeOSM(pt)
	return nil
}
