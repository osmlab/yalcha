package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/osmlab/gomap/gomap"
)

// Server contains Openstreetmap API handlers
type Server struct {
	g *gomap.Gomap
}

// SetHeaders is used to set default headers for OK response
func (s *Server) SetHeaders(c echo.Context) {
	c.Response().Header().Set(echo.HeaderContentType, strings.ToLower(echo.MIMETextXMLCharsetUTF8))
	c.Response().Header().Set("Cache-Control", "private, max-age=0, must-revalidate")
	c.Response().WriteHeader(http.StatusOK)
}

// SetEmptyResultHeaders is used to set specific headers for empty result
func (s *Server) SetEmptyResultHeaders(c echo.Context, status int) {
	c.Response().Header().Set(echo.HeaderContentType, strings.ToLower(echo.MIMETextXMLCharsetUTF8))
	c.Response().Header().Set(echo.HeaderContentLength, "0")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().WriteHeader(status)
}

// New returns new Server
func New(g *gomap.Gomap) *Server {
	return &Server{g: g}
}
