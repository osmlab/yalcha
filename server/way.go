package server

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/osmlab/gomap/osm"
)

// GetWay returns way by id
func (s *Server) GetWay(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	wayID, err := s.db.GetWayID(id)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	isVisible, err := s.db.IsWayVisible(wayID)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}
	if !isVisible {
		s.SetEmptyResultHeaders(c, http.StatusGone)
		return nil
	}

	ids := []int64{wayID}
	ways, err := s.db.GetWays(ids)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}

	resp := osm.New()
	resp.Ways = ways

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}

// GetWayByVersion returns way by id and version
func (s *Server) GetWayByVersion(c echo.Context) error {
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
	way, err := s.db.GetWayByVersion(id, version)
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

	wayID, err := s.db.GetWayID(id)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	isVisible, err := s.db.IsWayVisible(wayID)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}
	if !isVisible {
		s.SetEmptyResultHeaders(c, http.StatusGone)
		return nil
	}

	ids := []int64{wayID}
	ways, err := s.db.GetWays(ids)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}

	wayIDs := []int64{}
	for i := range ways {
		wayIDs = append(wayIDs, ways[i].ID)
	}
	nodeIDs, err := s.db.GetNodesFromWays(wayIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}
	nodes, err := s.db.GetNodes(nodeIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}

	resp := osm.New()
	resp.Nodes = nodes
	resp.Ways = ways

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}

// GetWayHistory returns way history by id
func (s *Server) GetWayHistory(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	ways, err := s.db.GetWayHistory(id)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	resp := osm.New()
	resp.Ways = ways

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

	wayIDs := make([]int64, 0)
	historicWayIDs := make([][2]int64, 0)
	for i := range wayIDsString {
		idv := strings.Split(wayIDsString[i], "v")
		id, err := strconv.ParseInt(idv[0], 10, 64)
		if err != nil {
			s.SetEmptyResultHeaders(c, http.StatusBadRequest)
			return nil
		}
		if len(idv) == 1 {
			wayIDs = appendIfUnique(wayIDs, id)
			continue
		}
		v, err := strconv.ParseInt(idv[1], 10, 64)
		if err != nil {
			s.SetEmptyResultHeaders(c, http.StatusBadRequest)
			return nil
		}
		historicWayIDs = appendVersionIfUnique(historicWayIDs, [2]int64{id, v})
	}

	currentWays, err := s.db.GetWays(wayIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	historicWays, err := s.db.GetHistoricWays(historicWayIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	ways := append(currentWays, historicWays...)

	if len(ways) != len(wayIDs)+len(historicWayIDs) {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return nil
	}

	resp := osm.New()
	resp.Ways = ways

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}
