package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type CommonConfiguration struct {
	ServicePort       int    `default:"3000"`
	ReadTimeoutInSec  int    `default:"15"`
	WriteTimeoutInSec int    `default:"15"`
	IdleTimeoutInSec  int    `deault:"60"`
	StoreType         string `default:"local"`
}

type LogConfiguration struct {
	Level    string `default:"-1"`
	Format   string `default:"_2 Jan 2006 15:04:05.000"`
	Encoding string `default:"console"`
}

var (
	CommonConfig CommonConfiguration
	LogConfig    LogConfiguration
)

func LoadConfig() error {
	// loading common config
	if err := envconfig.Process("", &CommonConfig); err != nil {
		log.Printf("Failed to load common config env %v\n", err)
		return err
	}
	log.Printf("CommonConfig: %+v\n", CommonConfig)

	// loading log config
	if err := envconfig.Process("", &LogConfig); err != nil {
		log.Printf("Failed to load common log env %v\n", err)
		return err
	}
	log.Printf("LogConfig: %+v\n", LogConfig)
	return nil
}
