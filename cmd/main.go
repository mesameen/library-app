package main

import (
	"github.com/test/library-app/internal/config"
	"github.com/test/library-app/internal/logger"
)

func main() {
	// loads config if any error in reading config panics the appl
	config.LoadConfig()

	// configures logger for an app
	logger.InitLogger()
	logger.Infof("Hello this is library-app")
}
