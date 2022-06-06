package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bjornaer/sailor/internal/port"
	"github.com/gin-gonic/gin"
)

func digestPorts(c *gin.Context) {
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

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the Port Domain Service!")
	})
	router.GET("/digest", digestPorts)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
