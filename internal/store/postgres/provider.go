package postgres

import (
	"context"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/test/library-app/internal/config"
	"github.com/test/library-app/internal/logger"
)

func InitPostgresStore() (*PostgresDB, error) {
	query := url.Values{}
	query.Add("application_name", config.CommonConfig.AppName)
	query.Add("client_encoding", "utf-8")
	query.Add("sslmode", "disable")
	u := url.URL{
		Scheme:   "postgres",
		Host:     config.PostgresConfig.Host,
		User:     url.UserPassword(config.PostgresConfig.PGUserName, config.PostgresConfig.Password),
		Path:     "/" + config.PostgresConfig.DBName,
		RawQuery: query.Encode(),
	}
	logger.Infof("Connecting to postgres: %s", u.String())
	poolConfig, err := pgxpool.ParseConfig(u.String())
	if err != nil {
		logger.Errorf("Failed to parse postgres config. Error: %v", err)
		return nil, err
	}
	poolConfig.MaxConns = 10
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Errorf("Failed to create postgres connection pool. Error: %v", err)
		return nil, err
	}
	// checking by pinging to DB
	err = pool.Ping(ctx)
	if err != nil {
		logger.Errorf("Failed to ping to connected postgres. Error: %v", err)
		return nil, err
	}
	logger.Infof("Connected to postgress successfully")
	return &PostgresDB{
		DB: pool,
	}, nil
}
