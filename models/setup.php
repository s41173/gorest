package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// debug env (hapus kalau sudah jalan)
	// log.Println("MYSQLUSER:", os.Getenv("MYSQLUSER"))
	// log.Println("MYSQLHOST:", os.Getenv("MYSQLHOST"))
	// log.Println("MYSQLPORT:", os.Getenv("MYSQLPORT"))
	// log.Println("MYSQLDATABASE:", os.Getenv("MYSQLDATABASE"))

	// buat DSN dari environment variable
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQLUSER"),
		os.Getenv("MYSQLPASSWORD"),
		os.Getenv("MYSQLHOST"),
		os.Getenv("MYSQLPORT"),
		os.Getenv("MYSQLDATABASE"),
	)

	// log.Println("Connecting to DB:", dsn)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect database: %v", err)
	}

	log.Println("Database connected successfully")

	DB = database
}
