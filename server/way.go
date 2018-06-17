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
	osm, err := s.db.GetWayFull(id)
	if err != nil || len(osm.Objects()) == 0 {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	if !osm.Ways[0].Visible {
		s.SetEmptyResultHeaders(c, http.StatusGone)
		return nil
	}

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(osm)
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

	cIDs := make([]int64, 0)
	nIDsVs := make([][2]int64, 0)
	for i := range wayIDsString {
		idv := strings.Split(wayIDsString[i], "v")
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

	ways, err := s.db.GetWays(cIDs, nIDsVs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	if len(ways) != len(cIDs)+len(nIDsVs) {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return nil
	}

	resp := osm.New()
	resp.Ways = ways

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}
