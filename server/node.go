package server

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/osmlab/gomap/gomap"
)

// GetNode returns node by id
func (s *Server) GetNode(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	resp, err := s.g.NodeHandler(id)
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

// GetNodes returns nodes by ids
func (s *Server) GetNodes(c echo.Context) error {
	rawIDs := strings.Split(c.QueryParam("nodes"), ",")
	if len(rawIDs) == 0 {
		s.SetEmptyResultHeaders(c, http.StatusBadRequest)
		return nil
	}

	ids, histIDs, err := getCurrentHistoricIDs(rawIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusBadRequest)
		return err
	}

	resp, err := s.g.NodesHandler(ids, histIDs)
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
