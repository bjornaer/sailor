package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bjornaer/sailor/internal/port"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello Sailor! Welcome to the Port Domain Service!")
	})
	router.GET("/process", processPorts)
	return router
}

func processPorts(c *gin.Context) {
	fileName := os.Getenv("PORTS_FILE")
	if len(fileName) == 0 {
		fileName = "./ports.json"
	}
	err := port.HandlePorts(fileName, func(portKey string, portData port.PortData) {
		fmt.Printf("Port Key %s : Port Data %v", portKey, portData.Name)
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.String(http.StatusOK, "Finished updating DB with ports data!")
}