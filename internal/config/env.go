package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST     string
	DB_PASSWORD string
	DB_USER     string
	DB_NAME     string
	DB_PORT     int
}

func New() *Config {
	return &Config{}
}

func (c *Config) Load() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	c.DB_HOST = os.Getenv("DB_HOST")
	c.DB_NAME = os.Getenv("DB_NAME")
	c.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	c.DB_USER = os.Getenv("DB_USER")
	c.DB_PORT, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}

	return nil
}
