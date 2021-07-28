package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver string `mapstructure:"db_drive"`
	DBSource string `mapstructure:"db_source"`
}

func createConnection() (*sql.DB, error) {

	config, cfErr := loadConfig(".")

	if cfErr != nil {
		return nil, cfErr
	}

	// open database
	db, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		return nil, err
	}

	// check db
	err = db.Ping()

	if err != nil {
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
