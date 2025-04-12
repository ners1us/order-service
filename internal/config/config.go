package config

import "os"

type Config struct {
	DbUrl     string
	JWTSecret string
	RestPort  string
}

func NewConfig() *Config {
	return &Config{
		DbUrl:     getEnv("DB_URL"),
		JWTSecret: getEnv("JWT_SECRET"),
		RestPort:  getEnv("REST_PORT"),
	}
}

func getEnv(key string) string {
	value, _ := os.LookupEnv(key)
	return value
}
