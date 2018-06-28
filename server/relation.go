package server

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/osmlab/gomap/gomap"
)

// GetRelation returns relation by id
func (s *Server) GetRelation(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	resp, err := s.g.RelationHandler(id)
	if err == gomap.ErrElementNotFound {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	if err == gomap.ErrElementDeleted {
		s.SetEmptyResultHeaders(c, http.StatusGone)
		return err
	}
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}

// GetRelations returns relations by ids
func (s *Server) GetRelations(c echo.Context) error {
	rawIDs := strings.Split(c.QueryParam("relations"), ",")
	if len(rawIDs) == 0 {
		s.SetEmptyResultHeaders(c, http.StatusBadRequest)
		return nil
	}

	ids, histIDs, err := getCurrentHistoricIDs(rawIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusBadRequest)
		return err
	}

	resp, err := s.g.RelationsHandler(ids, histIDs)
	if err == gomap.ErrElementNotFound {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}

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

	resp, err := s.g.RelationFullHandler(id)
	if err == gomap.ErrElementNotFound {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	if err == gomap.ErrElementDeleted {
		s.SetEmptyResultHeaders(c, http.StatusGone)
		return err
	}
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}

// GetRelationByVersion returns relation by id and version
func (s *Server) GetRelationByVersion(c echo.Context) error {
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

	resp, err := s.g.NodeVersionHandler(id, version)
	if err == gomap.ErrElementNotFound {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusInternalServerError)
		return err
	}

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
}
