package main

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver string `mapstructure:"db_driver"`
	DBSource string `mapstructure:"db_source"`
}

func createConnection(cfg *Config) (*sql.DB, error) {
	// open database
	db, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		return nil, err
	}

	// check db
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func loadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("dbconfig")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
