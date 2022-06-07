package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bjornaer/sailor/internal/db"
	s "github.com/bjornaer/sailor/internal/sessionmanager"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if len(redisAddr) == 0 {
		redisAddr = "localhost:6379"
	}
	portsFileName := os.Getenv("PORTS_FILE")
	if len(portsFileName) == 0 {
		portsFileName = "./ports.json"
	}
	dbClient, err := db.InitDBClient() // "redis", redisAddr
	if err != nil {
		log.Fatal("Failed initialising DB connection")
	}
	sm := &s.SessionManager{UpdatesFile: portsFileName, DBClient: dbClient}
	router := SetupRouter(sm)

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
