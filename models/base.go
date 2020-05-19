package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var db *gorm.DB // database
var dbURI string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}

	dbURI = os.Getenv("DB_URI")

	migrateDatabase()
}

// Datebase migration
func migrateDatabase() {
	db := GetDB()

	db.Debug().AutoMigrate(&Movie{})
}

func GetDB() *gorm.DB {
	// Making connection to the database
	db, err := gorm.Open("mysql", dbURI)
	if err != nil {
		log.Println(err)
	}

	retryCount := 30
	for {
		err := db.DB().Ping()
		if err != nil {
			if retryCount == 0 {
				log.Fatalf("Not able to establish connection to database")
			}

			log.Printf(fmt.Sprintf("Could not connect to database. Wait 2 seconds. %d retries left...", retryCount))
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	return db
}
