package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bjornaer/sailor/internal/db"
	"github.com/bjornaer/sailor/internal/port"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello Sailor! Welcome to the Port Domain Service!")
	})
	router.GET("/process", processPorts)
	return router
}

func processPorts(c *gin.Context) {
	etcdAddr := os.Getenv("ETCD_ADDR")
	if len(etcdAddr) == 0 {
		etcdAddr = "127.0.0.1:2379"
	}
	db_cli := db.NewEtcdClient(etcdAddr) // TODO create a session manager that holds the DB connection and configs
	fileName := os.Getenv("PORTS_FILE")
	if len(fileName) == 0 {
		fileName = "./ports.json"
	}
	err := port.HandlePorts(fileName, func(portKey string, portData port.PortData) {
		fmt.Printf("| Port Key: %s | Port Name: %v | \n", portKey, portData.Name)
		portDataStr, err := portData.Serialize()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		db_cli.KV.Put(db_cli.Ctx, portKey, portDataStr)
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.String(http.StatusOK, "Finished updating DB with ports data!")
}
