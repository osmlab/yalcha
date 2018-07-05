package server

import (
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/osmlab/gomap/gomap"
)

// GetMap returns amp elements
func (s *Server) GetMap(c echo.Context) error {
	bboxRaw := strings.Split(c.Param("bbox"), ",")
	if len(bboxRaw) != 3 {
		return errors.New("malo argumentov")
	}

	bbox := []float64{}
	for i := range bboxRaw {
		arg, err := strconv.ParseFloat(bboxRaw[i], 64)
		if err != nil {
			return err
		}
		bbox = append(bbox, arg)
	}

	resp, err := s.g.MapHandler(bbox)
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
