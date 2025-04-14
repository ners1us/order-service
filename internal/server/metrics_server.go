package server

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

type metricsServer struct {
	server *http.Server
	mux    *http.ServeMux
}

func NewMetricsServer(port string) BackendServer {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return &metricsServer{
		server: srv,
		mux:    mux,
	}
}

func (ms *metricsServer) ConfigureRoutes() {
	ms.mux.Handle("/metrics", promhttp.Handler())
}

func (ms *metricsServer) Start() error {
	log.Println("starting metrics server...")
	return ms.server.ListenAndServe()
}

func (ms *metricsServer) Stop(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		err := ms.server.Shutdown(shutdownCtx)
		errCh <- err
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Printf("metrics server shutdown error: %v", err)
		} else {
			log.Println("metrics server stopped gracefully")
		}
	case <-shutdownCtx.Done():
		log.Println("metrics server forced to stop")
	}
}
