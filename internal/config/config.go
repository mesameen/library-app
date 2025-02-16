package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type CommonConfiguration struct {
	AppName           string `default:"library-app"`
	ServicePort       int    `default:"3000"`
	ReadTimeoutInSec  int    `default:"15"`
	WriteTimeoutInSec int    `default:"15"`
	IdleTimeoutInSec  int    `deault:"60"`
	StoreType         string `default:"local"` // local | postgres
}

type LogConfiguration struct {
	Level    string `default:"-1"`
	Format   string `default:"_2 Jan 2006 15:04:05.000"`
	Encoding string `default:"console"`
}

type PostgresConfiguration struct {
	Host           string `default:"localhost:5432"`
	PGUserName     string `default:"postgres"`
	Password       string `default:"postgres"`
	DBName         string `default:"postgresdb"`
	BooksTableName string `default:"books"`
	LoansTableName string `default:"loans"`
}

var (
	CommonConfig   CommonConfiguration
	LogConfig      LogConfiguration
	PostgresConfig PostgresConfiguration
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

	// loading postgres config
	if err := envconfig.Process("", &PostgresConfig); err != nil {
		log.Printf("Failed to load common log env %v\n", err)
		return err
	}
	log.Printf("PostgresConfig: %+v\n", PostgresConfig)

	return nil
}
