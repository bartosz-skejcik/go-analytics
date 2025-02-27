package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bartosz-skejcik/go-analytics/internal/config"
	"github.com/bartosz-skejcik/go-analytics/internal/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Db     *sql.DB
	Client *gorm.DB
	Config *config.Config
}

func New(c *config.Config) *Database {
	return &Database{
		Config: c,
	}
}

func (d *Database) Connect() error {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", d.Config.DB_HOST, d.Config.DB_USER, d.Config.DB_PASSWORD, d.Config.DB_NAME, d.Config.DB_PORT)

	d.Client, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	CheckError(err)

	return err
}

func (d *Database) RunMigrations() {
	err := d.Connect()
	if err != nil {
		panic(err)
	}

	d.Db, err = d.Client.DB()
	if err != nil {
		panic(err)
	}

	defer d.Db.Close()

	err = d.Client.AutoMigrate(&models.Session{}, &models.PageView{}, &models.Event{})
	if err != nil {
		log.Printf("Error while migrating: %v", err)
		return
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
