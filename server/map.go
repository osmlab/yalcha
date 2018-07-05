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
	bboxRaw := strings.Split(c.QueryParam("bbox"), ",")
	if len(bboxRaw) != 4 {
		return errors.New("few arguments")
	}

	bbox := []int64{}
	for i := range bboxRaw {
		bboxRaw[i] += "0000000000"
		arg, err := strconv.ParseInt(strings.Replace(bboxRaw[i], ".", "", -1)[:8], 10, 64)
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
