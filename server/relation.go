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
	relation, err := s.db.GetRelation(id)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	if !relation.Visible {
		s.SetEmptyResultHeaders(c, http.StatusGone)
		return nil
	}

	relationIDs := []int64{}
	nodeIDs := []int64{}
	wayIDs := []int64{}
	for _, m := range relation.Members {
		if m.Type == "relation" {
			relationIDs = append(relationIDs, m.Ref)
		}
		if m.Type == "way" {
			wayIDs = append(wayIDs, m.Ref)
		}
		if m.Type == "node" {
			nodeIDs = append(nodeIDs, m.Ref)
		}
	}

	relations, err := s.db.GetRelations(relationIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	ways, err := s.db.GetWays(wayIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	for _, w := range *ways {
		for _, wn := range w.Nodes {
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
	}
	nodes, err := s.db.GetNodes(nodeIDs)
	if err != nil {
		s.SetEmptyResultHeaders(c, http.StatusNotFound)
		return err
	}

	resp := osm.New()
	resp.Nodes = *nodes
	resp.Ways = *ways
	resp.Relations = *relations
	resp.Relations = append(resp.Relations, relation)

	s.SetHeaders(c)
	return xml.NewEncoder(c.Response()).Encode(resp)
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
