package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/osmlab/gomap/config"
	"github.com/osmlab/gomap/server"
)

// Load returns api router
func Load(config *config.Config, s *server.Server) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{}))

	api := e.Group("/api")

	api06 := api.Group("/0.6")

	node06 := api06.Group("/node")
	node06.HEAD("/:id", s.GetNode)
	node06.GET("/:id", s.GetNode)
	node06.HEAD("/:id/:version", s.GetNodeByVersion)
	node06.GET("/:id/:version", s.GetNodeByVersion)
	node06.HEAD("/:id/history", s.GetNodeHistory)
	node06.GET("/:id/history", s.GetNodeHistory)
	node06.HEAD("/:id/ways", s.GetWaysForNode)
	node06.GET("/:id/ways", s.GetWaysForNode)

	nodes06 := api06.Group("/nodes")
	nodes06.HEAD("", s.GetNodes)
	nodes06.GET("", s.GetNodes)

	way06 := api06.Group("/way")
	way06.HEAD("/:id", s.GetWay)
	way06.GET("/:id", s.GetWay)
	way06.HEAD("/:id/:version", s.GetWayByVersion)
	way06.GET("/:id/:version", s.GetWayByVersion)
	way06.HEAD("/:id/full", s.GetWayFull)
	way06.GET("/:id/full", s.GetWayFull)
	way06.HEAD("/:id/history", s.GetWayHistory)
	way06.GET("/:id/history", s.GetWayHistory)

	ways06 := api06.Group("/ways")
	ways06.HEAD("", s.GetWays)
	ways06.GET("", s.GetWays)

	relation06 := api06.Group("/relation")
	relation06.HEAD("/:id", s.GetRelation)
	relation06.GET("/:id", s.GetRelation)
	relation06.HEAD("/:id/:version", s.GetRelationByVersion)
	relation06.GET("/:id/:version", s.GetRelationByVersion)
	relation06.HEAD("/:id/full", s.GetRelationFull)
	relation06.GET("/:id/full", s.GetRelationFull)
	relation06.HEAD("/:id/history", s.GetRelationHistory)
	relation06.GET("/:id/history", s.GetRelationHistory)

	relations06 := api06.Group("/relations")
	relations06.HEAD("", s.GetRelations)
	relations06.GET("", s.GetRelations)

	return e
}
