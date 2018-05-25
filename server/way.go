package server

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/osmlab/yalcha/osm"
)

// GetWay returns way by id
func (s *Server) GetWay(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	way, err := s.db.GetWay(id)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	if !way.Visible {
		s.SetEmptyResultHeaders(c, http.StatusGone)
		return nil
	}

	resp := osm.New()
	resp.Ways = append(resp.Ways, way)

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}

// GetWayFull returns full way by id
func (s *Server) GetWayFull(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	way, err := s.db.GetWay(id)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	if !way.Visible {
		s.SetEmptyResultHeaders(c, http.StatusGone)
		return nil
	}

	nodeIDs := []int64{}
	for _, wn := range way.Nodes {
		exist := false
		for _, id := range nodeIDs {
			if wn.ID == id {
				exist = true
				break
			}
		}
		if !exist {
			nodeIDs = append(nodeIDs, wn.ID)
		}
	}
	nodes, err := s.db.GetNodes(nodeIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	resp := osm.New()
	resp.Nodes = *nodes
	resp.Ways = append(resp.Ways, way)

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}

// GetWays returns ways by ids
func (s *Server) GetWays(c echo.Context) error {
	wayIDsString := strings.Split(c.QueryParam("ways"), ",")
	if len(wayIDsString) == 0 {
		s.SetEmptyResultHeaders(c, http.StatusBadRequest)
		return nil
	}

	wayIDs := []int64{}
	for i := range wayIDsString {
		id, err := strconv.ParseInt(wayIDsString[i], 10, 64)
		if err != nil {
			s.SetEmptyResultHeaders(c, http.StatusBadRequest)
			return nil
		}
		wayIDs = append(wayIDs, id)
	}

	ways, err := s.db.GetWays(wayIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	if len(*ways) != len(wayIDs) {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return nil
	}

	resp := osm.New()
	resp.Ways = *ways

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}
