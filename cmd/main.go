package main

import (
	"os"
	"yubigo-pass/internal/app"
	"yubigo-pass/internal/app/services"

	"github.com/sirupsen/logrus"
)

func main() {
	err := app.NewRunner(services.Build()).Run()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
