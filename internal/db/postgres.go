package db

import (
	"fmt"
	"log"
	"os"

	"github.com/bartosz-skejcik/go-analytics/internal/config"
	"github.com/go-pg/pg/v10"
)

type Database struct {
	Db     *pg.DB
	config *config.Config
}

func New(c *config.Config) *Database {
	return &Database{
		Db:     nil,
		config: c,
	}
}

func (d *Database) Connect() error {
	var err error

	d.Db = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", d.config.DB_HOST, d.config.DB_PORT),
		User:     d.config.DB_USER,
		Password: d.config.DB_PASSWORD,
		Database: d.config.DB_NAME,
	})
	CheckError(err)

	return err
}

func (d *Database) RunMigrations() error {
	// 1. read the migrations folder for files
	// 2. read each file and save the contents to a string variable
	// 3. execute the string variable on the database

	var err error

	err = d.Connect()
	if err != nil {
		return err
	}

	defer d.Db.Close()

	migrationsDir := "migrations"

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", migrationsDir, file.Name())

		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		log.Printf("Executing %s", file.Name())

		_, err = d.Db.Exec(string(content))
		if err != nil {
			return err
		}
	}

	log.Println("Finished running all migrations")

	return nil
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
