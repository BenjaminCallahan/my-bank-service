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
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	handlers := api.NewHandler()

	srv := api.NewServer("", handlers.InitRoutes())
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
