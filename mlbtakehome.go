package main

import (
	"github.com/stefanKnott/mlbtakehome/pkg/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// start server
	handlers.InitTeamIdSet()
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/liveness", handlers.GetLiveness)
		v1.GET("/readiness", handlers.GetReadiness)
		v1.GET("/schedule", handlers.GetSchedule)
	}
	router.Run()
}
