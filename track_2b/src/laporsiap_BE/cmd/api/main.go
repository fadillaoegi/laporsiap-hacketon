package main

import (
	"context"
	"laporsiap/internal/application"
	"laporsiap/internal/infrastructure/apertus"
	"laporsiap/internal/infrastructure/config"
	"laporsiap/internal/infrastructure/postgres"
	"laporsiap/internal/infrastructure/reports"
	"laporsiap/internal/transport/httpapi"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	limits, err := config.LoadSafetyLimits()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	store, err := postgres.OpenStore(ctx, os.Getenv("DATABASE_URL"))
	cancel()
	if err != nil {
		log.Fatal("PostgreSQL connection/migration failed; check local configuration.")
	}
	defer store.Pool.Close()
	store.Limits = limits
	cleanupCtx, stopCleanup := context.WithCancel(context.Background())
	defer stopCleanup()
	go func() {
		for {
			ctx, cancel := context.WithTimeout(cleanupCtx, 10*time.Second)
			_, err := store.PurgeExpired(ctx)
			cancel()
			if err != nil && cleanupCtx.Err() == nil {
				log.Print("Expired-case cleanup failed; access remains blocked after expiry.")
			}
			select {
			case <-cleanupCtx.Done():
				return
			case <-time.After(time.Hour):
			}
		}
	}()
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8787"
	}
	host, port, err := net.SplitHostPort(addr)
	containerBind := host == "0.0.0.0" && os.Getenv("ALLOW_CONTAINER_BIND") == "true"
	if err != nil || (host != "127.0.0.1" && host != "localhost" && !containerBind) || port != "8787" {
		log.Fatal("Local MVP must listen on 127.0.0.1:8787 or localhost:8787.")
	}
	provider := apertus.NewProvider(apertus.ConfigFromEnv())
	if provider.Key != "" {
		// Keep serving: every generation re-validates and the manual path stays available.
		if err := provider.ValidateConfiguration(); err != nil {
			log.Printf("AI provider disabled: %v", err)
		}
	}
	handler := (&httpapi.Server{Service: application.NewService(store, provider, reports.Renderer{})}).Handler()
	if directory := os.Getenv("FRONTEND_DIR"); directory != "" {
		handler, err = httpapi.FrontendHandler(directory, handler)
		if err != nil {
			log.Fatal("Compiled frontend unavailable.")
		}
	}
	server := &http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 60 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16 * 1024}
	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopped
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()
	log.Printf("LaporSiap API ready at http://%s (PostgreSQL; local only)", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("API listener failed.")
	}
}
