package server

import (
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/osmlab/gomap/gomap"
)

// GetChangeset returns changeset by id
func (s *Server) GetChangeset(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	includeDiscussion := false
	includeDiscussionRaw := c.QueryParam("include_discussion")
	if len(includeDiscussionRaw) != 0 {
		includeDiscussion, err = strconv.ParseBool(includeDiscussionRaw)
		if err != nil {
			s.SetEmptyResultHeaders(c, http.StatusNotFound)
			return err
		}
	}

	resp, err := s.g.ChangesetHandler(id, includeDiscussion)
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
