package models

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	err := mysqlDriver.RegisterTLSConfig("custom", tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=custom&parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect database: %v", err)
	}

	DB = database
}
