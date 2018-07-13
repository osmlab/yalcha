package osm

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"sort"
	"strconv"
)

// These values should be returned if the osm data is actual
// osm data to give some information about the source and license.
const (
	Copyright   = "OpenStreetMap and contributors"
	Attribution = "http://www.openstreetmap.org/copyright"
	License     = "http://opendatacommons.org/licenses/odbl/1-0/"
	Version     = 0.6
	Generator   = "Gomap"
)

// OSM represents the core osm data
// designed to parse http://wiki.openstreetmap.org/wiki/OSM_XML
type OSM struct {
	Version   float64 `xml:"version,attr,omitempty"`
	Generator string  `xml:"generator,attr,omitempty"`

	// These three attributes are returned by the osm api.
	// The Copyright, Attribution and License constants contain
	// suggested values that match those returned by the official api.
	Copyright   string `xml:"copyright,attr,omitempty"`
	Attribution string `xml:"attribution,attr,omitempty"`
	License     string `xml:"license,attr,omitempty"`

	Nodes      Nodes      `xml:"node"`
	Ways       Ways       `xml:"way"`
	Relations  Relations  `xml:"relation"`
	Changesets Changesets `xml:"changeset"`
}

// New creates osm object
func New() *OSM {
	return &OSM{
		Version:     Version,
		Generator:   Generator,
		Copyright:   Copyright,
		Attribution: Attribution,
		License:     License,
	}
}

// Objects returns an array of objects containing any nodes, ways, relations,
// changesets, notes and users.
func (o *OSM) Objects() Objects {
	if o == nil {
		return nil
	}

	result := make(Objects, 0, len(o.Nodes)+len(o.Ways)+len(o.Relations)+len(o.Changesets))
	for _, o := range o.Nodes {
		result = append(result, o)
	}
	for _, o := range o.Ways {
		result = append(result, o)
	}
	for _, o := range o.Relations {
		result = append(result, o)
	}
	for _, o := range o.Changesets {
		result = append(result, o)
	}

	return result
}

// MarshalJSON allows the tags to be marshalled as an object
// as defined by the overpass osmjson.
// http://overpass-api.de/output_formats.html#json
func (o OSM) MarshalJSON() ([]byte, error) {
	s := struct {
		Version     float64 `json:"version,omitempty"`
		Generator   string  `json:"generator,omitempty"`
		Copyright   string  `json:"copyright,omitempty"`
		Attribution string  `json:"attribution,omitempty"`
		License     string  `json:"license,omitempty"`
		Elements    Objects `json:"elements"`
	}{o.Version, o.Generator, o.Copyright,
		o.Attribution, o.License, o.Objects()}

	return json.Marshal(s)
}

// MarshalXML implements the xml.Marshaller method to allow for the
// correct wrapper/start element case and attr data.
func (o OSM) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "osm"
	start.Attr = make([]xml.Attr, 0, 5)

	if o.Version != 0 {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: "version"},
			Value: strconv.FormatFloat(o.Version, 'g', -1, 64),
		})
	}

	if o.Generator != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "generator"}, Value: o.Generator})
	}

	if o.Copyright != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "copyright"}, Value: o.Copyright})
	}

	if o.Attribution != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "attribution"}, Value: o.Attribution})
	}

	if o.License != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "license"}, Value: o.License})
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := o.marshalInnerXML(e); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

func (o *OSM) marshalInnerXML(e *xml.Encoder) error {
	if o == nil {
		return nil
	}

	if err := e.Encode(o.Nodes); err != nil {
		return err
	}

	if err := e.Encode(o.Ways); err != nil {
		return err
	}

	if err := e.Encode(o.Relations); err != nil {
		return err
	}

	if err := e.Encode(o.Changesets); err != nil {
		return err
	}

	return nil
}

func (o *OSM) marshalInnerElementsXML(e *xml.Encoder) error {
	if err := e.Encode(o.Nodes); err != nil {
		return err
	}

	if err := e.Encode(o.Ways); err != nil {
		return err
	}

	if err := e.Encode(o.Relations); err != nil {
		return err
	}

	if err := e.Encode(o.Changesets); err != nil {
		return err
	}

	return nil
}

// Equals compares two OSM objects
func (o *OSM) Equals(obj *OSM) bool {
	o.sort()
	obj.sort()
	return reflect.DeepEqual(o.Objects(), obj.Objects())
}

func (o *OSM) sort() {
	for k := range o.Nodes {
		sort.Slice(o.Nodes[k].Tags, func(i, j int) bool {
			if o.Nodes[k].Tags[i].K < o.Nodes[k].Tags[j].K {
				return true
			}
			if o.Nodes[k].Tags[i].K > o.Nodes[k].Tags[j].K {
				return false
			}
			return o.Nodes[k].Tags[i].V < o.Nodes[k].Tags[j].V
		})
	}

	for k := range o.Ways {
		sort.Slice(o.Ways[k].Nodes, func(i, j int) bool {
			return o.Ways[k].Nodes[i].ID < o.Ways[k].Nodes[j].ID
		})
		sort.Slice(o.Ways[k].Tags, func(i, j int) bool {
			if o.Ways[k].Tags[i].K < o.Ways[k].Tags[j].K {
				return true
			}
			if o.Ways[k].Tags[i].K > o.Ways[k].Tags[j].K {
				return false
			}
			return o.Ways[k].Tags[i].V < o.Ways[k].Tags[j].V
		})
	}

	for k := range o.Relations {
		sort.Slice(o.Relations[k].Members, func(i, j int) bool {
			if o.Relations[k].Members[i].Ref < o.Relations[k].Members[j].Ref {
				return true
			}
			if o.Relations[k].Members[i].Ref > o.Relations[k].Members[j].Ref {
				return false
			}
			if o.Relations[k].Members[i].Type < o.Relations[k].Members[j].Type {
				return true
			}
			if o.Relations[k].Members[i].Type > o.Relations[k].Members[j].Type {
				return false
			}
			return o.Relations[k].Members[i].Role > o.Relations[k].Members[j].Role
		})
		sort.Slice(o.Relations[k].Tags, func(i, j int) bool {
			if o.Relations[k].Tags[i].K < o.Relations[k].Tags[j].K {
				return true
			}
			if o.Relations[k].Tags[i].K > o.Relations[k].Tags[j].K {
				return false
			}
			return o.Relations[k].Tags[i].V < o.Relations[k].Tags[j].V
		})
	}

	sort.Slice(o.Nodes, func(i, j int) bool { return o.Nodes[i].ID < o.Nodes[j].ID })
	sort.Slice(o.Ways, func(i, j int) bool { return o.Ways[i].ID < o.Ways[j].ID })
	sort.Slice(o.Relations, func(i, j int) bool { return o.Relations[i].ID < o.Relations[j].ID })
}
