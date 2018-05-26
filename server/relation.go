package server

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/osmlab/yalcha/osm"
)

// GetRelation returns relation by id
func (s *Server) GetRelation(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	relation, err := s.db.GetRelation(id)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	if !relation.Visible {
		s.SetEmptyResultHeaders(c, http.StatusGone)
		return nil
	}

	resp := osm.New()
	resp.Relations = append(resp.Relations, relation)

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}

// GetRelationFull returns full relation by id
func (s *Server) GetRelationFull(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	osm, err := s.db.GetRelationFull(id)
	if err != nil || len(osm.Objects()) == 0 {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(osm)
}

// GetRelations returns relations by ids
func (s *Server) GetRelations(c echo.Context) error {
	relationIDsString := strings.Split(c.QueryParam("relations"), ",")
	if len(relationIDsString) == 0 {
		s.SetEmptyResultHeaders(c, http.StatusBadRequest)
		return nil
	}

	relationIDs := []int64{}
	for i := range relationIDsString {
		id, err := strconv.ParseInt(relationIDsString[i], 10, 64)
		if err != nil {
			s.SetEmptyResultHeaders(c, http.StatusBadRequest)
			return nil
		}
		relationIDs = append(relationIDs, id)
	}

	relations, err := s.db.GetRelations(relationIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	if len(*relations) != len(relationIDs) {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return nil
	}

	resp := osm.New()
	resp.Relations = *relations

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}
