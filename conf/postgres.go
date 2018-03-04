package conf

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Sirupsen/logrus"
	"github.com/xlab/closer"

	_ "github.com/lib/pq" // Register pg driver
)

// PostgresConnect connects to postgres db
func PostgresConnect(config *Config) *sql.DB {

	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	closer.Bind(func() {
		logrus.Info("Closing database connection")
		db.Close()
	})

	return db
}
