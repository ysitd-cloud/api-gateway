package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"app.ysitd/gateway/pkg/bootstrap"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	grpcSrv := &http.Server{
		Addr:    ":50051",
		Handler: bootstrap.GetGrpcHandler(),
	}

	go func() {
		if err := grpcSrv.ListenAndServe(); err != nil {
			bootstrap.GetMainLogger().Error(err)
		}
	}()

	httpSrv := &http.Server{
		Addr: ":50050",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
				bootstrap.GetGrpcHandler().ServeHTTP(w, r)
			} else {
				bootstrap.GetHttpHandler().ServeHTTP(w, r)
			}
		}),
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			bootstrap.GetMainLogger().Error(err)
		}
	}()

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcSrv.Shutdown(ctx)
	httpSrv.Shutdown(ctx)
}
