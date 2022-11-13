package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type DBRepository struct {
	db *sqlx.DB
}

func (dbRepo *DBRepository) SetupDB() error {
	viper.SetConfigFile("config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.db-name"),
		viper.GetString("db.sslmode"))

	db, openErr := sqlx.Open("postgres", psqlInfo)
	if openErr != nil {
		return openErr
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	dbRepo.db = db
	return nil
}

func (dbRepo *DBRepository) Shutdown() error {
	return dbRepo.db.Close()
}
