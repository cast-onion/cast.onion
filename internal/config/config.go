package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	WebPort     string
	CDNPort     string
	DatabaseDSN string
	RedisAddr   string
	TLSCert     string
	TLSKey      string
	CDNBase     string
	AdminSecret string
	APIBase     string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("no .env file, using environment variables")
	}

	viper.SetDefault("PORT", "5000")
	viper.SetDefault("WEB_PORT", "2050")
	viper.SetDefault("CDN_PORT", "1050")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("CDN_BASE", "http://localhost:1050")
	viper.SetDefault("API_BASE", "http://localhost:5000")

	return &Config{
		Port:        viper.GetString("PORT"),
		WebPort:     viper.GetString("WEB_PORT"),
		CDNPort:     viper.GetString("CDN_PORT"),
		DatabaseDSN: viper.GetString("DATABASE_DSN"),
		RedisAddr:   viper.GetString("REDIS_ADDR"),
		TLSCert:     viper.GetString("TLS_CERT"),
		TLSKey:      viper.GetString("TLS_KEY"),
		CDNBase:     viper.GetString("CDN_BASE"),
		AdminSecret: viper.GetString("ADMIN_SECRET"),
		APIBase:     viper.GetString("API_BASE"),
	}
}
