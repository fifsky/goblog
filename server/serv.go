package server

import (
	"net/http"
	"log"
	"os"
	"os/signal"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

var serv = gin.Default()

func Serv() *gin.Engine {
	return serv
}

func Run(port string) {
	srv := &http.Server{
		Addr:    port,
		Handler: serv,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Close()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
