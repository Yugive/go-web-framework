package main

import (
	"context"
	"go-web-framework/framework"
	"go-web-framework/framework/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	core := framework.NewCore()
	core.Use(middleware.Recovery(), middleware.Cost())
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8084",
	}

	go func() {
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	timeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeCtx); err != nil {
		log.Fatal(err)
	}
}
