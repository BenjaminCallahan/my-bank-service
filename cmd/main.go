package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/BenjaminCallahan/my-bank-service/internal/api"
	"github.com/BenjaminCallahan/my-bank-service/internal/config"
	"github.com/BenjaminCallahan/my-bank-service/internal/repository"
	"github.com/BenjaminCallahan/my-bank-service/internal/repository/sqlite"
	"github.com/BenjaminCallahan/my-bank-service/internal/service"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	cfg, err := config.Init()
	if err != nil {
		logrus.Fatalf("failed to loading config: %s\n", err.Error())
	}

	db, err := sqlite.NewConnectDB(sqlite.Config{
		DBName: cfg.DBName,
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s\n", err.Error())
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handlers := api.NewHandler(service)

	srv := api.NewServer(cfg.Address, handlers.InitRoutes())
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("error occurred while running server: %s\n", err.Error())
		}
	}()

	logrus.Println("application started")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	logrus.Println("application shutdown")

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Printf("error occurred while shutting down http server: %s\n", err)
	}
}
