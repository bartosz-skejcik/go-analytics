package main

import (
	"github.com/bartosz-skejcik/go-analytics/internal/config"
	"github.com/bartosz-skejcik/go-analytics/internal/db"
	"github.com/bartosz-skejcik/go-analytics/internal/server"
)

func main() {
	conf := config.New()
	err := conf.Load()
	if err != nil {
		panic(err)
	}

	db := db.New(conf)

	db.RunMigrations()

	server := server.New(db)

	server.Run()
}
