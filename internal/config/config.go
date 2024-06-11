package config

import (
	"fmt"
	"flag"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var configPath string

func InitConfig(logger *zap.Logger) {
	flag.StringVar(&configPath, "config", "config.yaml", "Config file")
	flag.Parse()
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := ReadInConfig()
	if err != nil {
		logger.Panic("Fatal error config file", zap.Error(err))
	}
	viper.SetDefault("http_port", 8090)
	viper.SetDefault("pg_host", "postgres")
	viper.SetDefault("pg_port", 5432)
	viper.SetDefault("pg_user", "postgres")
	viper.SetDefault("pg_password", "password")
	viper.SetDefault("dbName", "deals")
	viper.SetDefault("autoMigrate", true)
	viper.SetDefault("enable_tracing", true)
	// viper.SetDefault("KEYCLOAK_URL", "http://localhost:8080")
	viper.SetDefault("KEYCLOAK_URL", "http://keycloak:8080")
	viper.SetDefault("KEYCLOAK_CLIENT_ID", "deals")
	viper.SetDefault("KEYCLOAK_CLIENT_SECRET", "5UOAC7zxygFkG6Qr6t1dsb5vOTpOHWTv")
	viper.SetDefault("KEYCLOAK_REALM", "master")
	viper.SetDefault("KEYCLOAK_DEFAULT_ROLE_ID", "ab6ccaee-6cba-4af8-b22b-96194097a0b4")
	viper.SetDefault("KEYCLOAK_DEFAULT_ROLE_NAME", "user")
	fmt.Println("KEYCLOAK_URL:", viper.GetString("KEYCLOAK_URL"))	
}

var ReadInConfig = func() error {
	return viper.ReadInConfig()
}

