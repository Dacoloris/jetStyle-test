package app

import (
	"flag"
	"fmt"
	"jetStyle-test/internal/config"
	"jetStyle-test/internal/handler"
	"jetStyle-test/internal/router"
	"jetStyle-test/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

type App struct {
	logger     *zap.Logger
	HTTPServer *http.Server
}

func New(cfg config.Config) (*App, error) {
	lg := logger.New()

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.Parse()

	h := handler.New(lg)

	srv := &http.Server{
		Handler: router.New(h),
		Addr:    fmt.Sprint(cfg.ServerAddress),
	}

	return &App{
		HTTPServer: srv,
		logger:     lg,
	}, nil
}
