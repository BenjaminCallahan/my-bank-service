package main

import (
	"github.com/sirupsen/logrus"

	"github.com/BenjaminCallahan/my-bank-service/internal/api"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	handlers := api.NewHandler()

	srv := api.NewServer("", handlers.InitRoutes())
	if err := srv.Run(); err != nil {
		logrus.Fatalf("error occurred while running server: %s\n", err.Error())
	}
}