package main

import (
	"context"
	"flag"
	url_shortener "github.com/Alzoww/url-shortener"
	"github.com/Alzoww/url-shortener/config"
	"github.com/Alzoww/url-shortener/internal/handler"
	"github.com/Alzoww/url-shortener/internal/storage/sqlite"
	"github.com/Alzoww/url-shortener/pkg/logger"
	"github.com/Alzoww/url-shortener/pkg/logger/sl"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var ushConfig config.Config
	var ushConfigPath string

	flag.StringVar(&ushConfigPath, "ushConfig", "config/local.yaml", "Path to ushConfig file")
	flag.Parse()

	err := config.LoadConfig(ushConfigPath, &ushConfig)
	if err != nil {
		panic(err)
	}

	log := logger.SetupLogger(ushConfig.Env)

	log.Info("starting url-shortener", slog.String("env", ushConfig.Env))
	log.Debug("debug messages are enabled")

	// sqlite storage
	storage, err := sqlite.New(ushConfig.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		return
	}

	h := handler.New(storage)

	server := new(url_shortener.Server)

	go func() {
		if err = server.Run(ushConfig.HttpServer, h.InitRoutes()); err != nil {
			log.Error("failed to run http server", sl.Err(err))
			os.Exit(1)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	if err = server.Shutdown(context.Background()); err != nil {
		log.Error("failed to shutdown server", sl.Err(err))
		os.Exit(1)
	}

	if err = storage.Close(); err != nil {
		log.Error("failed to close storage", sl.Err(err))
		os.Exit(1)
	}
}
