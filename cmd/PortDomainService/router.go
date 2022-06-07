package main

import (
	"net/http"

	s "github.com/bjornaer/sailor/internal/sessionmanager"
	"github.com/gin-gonic/gin"
)

func SetupRouter(sm *s.SessionManager) *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello Sailor! Welcome to the Port Domain Service!")
	})
	router.GET("/process", sm.ProcessPorts)
	router.GET("/port/:portid", sm.GetPort)
	return router
}
