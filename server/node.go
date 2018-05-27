package server

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/osmlab/yalcha/osm"
)

// GetNode returns node by id
func (s *Server) GetNode(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	node, err := s.db.GetNode(id)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	if !node.Visible {
		s.SetEmptyResultHeaders(c, http.StatusGone)
		return nil
	}

	resp := osm.New()
	resp.Nodes = append(resp.Nodes, node)

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}

// GetNodes returns nodes by ids
func (s *Server) GetNodes(c echo.Context) error {
	nodeIDsString := strings.Split(c.QueryParam("nodes"), ",")
	if len(nodeIDsString) == 0 {
		s.SetEmptyResultHeaders(c, http.StatusBadRequest)
		return nil
	}

	nodeIDs := []int64{}
	for i := range nodeIDsString {
		id, err := strconv.ParseInt(nodeIDsString[i], 10, 64)
		if err != nil {
			s.SetEmptyResultHeaders(c, http.StatusBadRequest)
			return nil
		}
		nodeIDs = append(nodeIDs, id)
	}

	nodes, err := s.db.GetNodes(nodeIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	if len(*nodes) != len(nodeIDs) {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return nil
	}

	resp := osm.New()
	for _, node := range *nodes {
		if !node.Visible {
			node.Lat = nil
			node.Lon = nil
		}
		resp.Nodes = append(resp.Nodes, node)
	}

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}
