package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type CommonConfiguration struct{}
type LogConfiguration struct {
	Level    string `default:"-1"`
	Format   string `default:"_2 Jan 2006 15:04:05.000"`
	Encoding string `default:"console"`
}

var (
	CommonConfig CommonConfiguration
	LogConfig    LogConfiguration
)

func LoadConfig() {
	// loading common config
	if err := envconfig.Process("", &CommonConfig); err != nil {
		log.Panicf("Failed to load common config env %v", err)
	}
	log.Printf("CommonConfig: %+v\n", CommonConfig)

	// loading log config
	if err := envconfig.Process("", &LogConfig); err != nil {
		log.Panicf("Failed to load common log env %v", err)
	}
	log.Printf("LogConfig: %+v\n", LogConfig)
}
