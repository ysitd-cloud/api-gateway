package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app.ysitd/gateway/pkg/bootstrap"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	server := &http.Server{
		Addr:    ":50051",
		Handler: bootstrap.GetHandler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			bootstrap.GetMainLogger().Error(err)
		}
	}()

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}
