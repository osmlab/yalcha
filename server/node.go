package server

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/osmlab/gomap/osm"
)

// GetNode returns node by id
func (s *Server) GetNode(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	nodeID, err := s.db.GetNodeID(id)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	isVisible, err := s.db.IsNodeVisible(nodeID)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}
	if !isVisible {
		s.SetEmptyResultHeaders(c, http.StatusGone)
		return nil
	}
	node, err := s.db.GetNode(nodeID)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}

	resp := osm.New()
	resp.Nodes = append(resp.Nodes, node)

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}

// GetNodeByVersion returns node by id and version
func (s *Server) GetNodeByVersion(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	version, err := strconv.ParseInt(c.Param("version"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	node, err := s.db.GetNodeByVersion(id, version)
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

// GetNodeHistory returns node history by id
func (s *Server) GetNodeHistory(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	nodes, err := s.db.GetNodeHistory(id)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	resp := osm.New()
	resp.Nodes = nodes

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}

// GetWaysForNode returns all the (not deleted) ways in which the given node is used
func (s *Server) GetWaysForNode(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	ways, err := s.db.GetWaysForNode(id)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	resp := osm.New()
	resp.Ways = *ways

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

	cIDs := make([]int64, 0)
	nIDsVs := make([][2]int64, 0)
	for i := range nodeIDsString {
		idv := strings.Split(nodeIDsString[i], "v")
		id, err := strconv.ParseInt(idv[0], 10, 64)
		if err != nil {
			s.SetEmptyResultHeaders(c, http.StatusBadRequest)
			return nil
		}
		if len(idv) == 1 {
			cIDs = appendIfUnique(cIDs, id)
			continue
		}
		v, err := strconv.ParseInt(idv[1], 10, 64)
		if err != nil {
			s.SetEmptyResultHeaders(c, http.StatusBadRequest)
			return nil
		}
		nIDsVs = appendVersionIfUnique(nIDsVs, [2]int64{id, v})
	}

	nodes, err := s.db.GetNodes(cIDs, nIDsVs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	if len(nodes) != len(cIDs)+len(nIDsVs) {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return nil
	}

	resp := osm.New()
	for _, node := range nodes {
		if !node.Visible {
			node.Lat = nil
			node.Lon = nil
		}
		resp.Nodes = append(resp.Nodes, node)
	}

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}
