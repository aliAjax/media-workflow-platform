package main

import (
	"context"
	"errors"
	"example.com/media-workflow/internal/application"
	"example.com/media-workflow/internal/observability"
	"example.com/media-workflow/internal/repository"
	"example.com/media-workflow/internal/scheduler"
	"example.com/media-workflow/internal/storage"
	"example.com/media-workflow/internal/transport/httpapi"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	l := observability.NewLogger()
	store, e := repository.New(env("MEDIA_STATE", "./data/state.json"))
	if e != nil {
		l.Error("store init", slog.Any("error", e))
		os.Exit(1)
	}
	r := scheduler.New(store, store, store, store, storage.NewLocalBlobStore(env("MEDIA_BLOBS", "./data/blobs")), l, 2)
	r.Start()
	svc := application.New(store, store, store, store, r.Enqueue)
	srv := &http.Server{Addr: env("MEDIA_ADDR", ":8084"), Handler: httpapi.New(svc, store, store, store).Routes(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		c, cc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cc()
		_ = srv.Shutdown(c)
		r.Stop()
	}()
	l.Info("listening", "addr", srv.Addr)
	if e := srv.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
		l.Error("server", slog.Any("error", e))
		os.Exit(1)
	}
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
