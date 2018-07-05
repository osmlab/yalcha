package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/osmlab/gomap/osm"
)

const (
	gomapBaseURL  = "http://localhost:8090"
	cgimapBaseURL = "http://localhost:31337"
)

var routes = []string{
	"/api/0.6/node/1001",
	"/api/0.6/node/1001/history",
	"/api/0.6/node/1001/2",
	"/api/0.6/node/1004/ways",
	"/api/0.6/nodes?nodes=1001,1002,1003,1005v1",

	"/api/0.6/way/3001",
	"/api/0.6/way/3004/full",
	"/api/0.6/way/3004/history",
	"/api/0.6/way/3004/1",
	"/api/0.6/ways?ways=3001,3004,3006",

	"/api/0.6/relation/8005",
	"/api/0.6/relation/5006/full",
	"/api/0.6/relation/8005/history",
	"/api/0.6/relation/8005/1",
	"/api/0.6/relations?relations=8001,8005v1",

	"/api/0.6/map?bbox=1.0010000,1.0010000,1.0060000,1.7030000",
}

func main() {
	for _, r := range routes {
		gomap := makeRequest(gomapBaseURL + r)
		cgimap := makeRequest(cgimapBaseURL + r)
		isEqual := gomap.Equals(cgimap)
		log.Printf("%v: %v", r, isEqual)
		if !isEqual {
			log.Println(gomap)
			log.Println(cgimap)
		}
	}
}

func makeRequest(url string) *osm.OSM {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	o := &osm.OSM{}
	xml.Unmarshal(response, o)
	return o
}
