package database

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DatabaseConnection = func() *gorm.DB {
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", viper.GetString("pg_host"), viper.GetInt("pg_port"), viper.GetString("pg_user"), viper.GetString("pg_password"), viper.GetString("dbName"))
	db, err := Open(postgres.Open(sqlInfo))
	if err != nil {
		panic(err)
	}

	return db
}

var Open = func(dialector gorm.Dialector) (*gorm.DB, error) {
	return gorm.Open(dialector, &gorm.Config{})
}
