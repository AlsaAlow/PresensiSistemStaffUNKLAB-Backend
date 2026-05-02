package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	// 🔥 WAJIB: parseTime=true supaya DATETIME bisa dibaca sebagai time.Time
	dsn := "root:123456@tcp(127.0.0.1:3306)/sipres?parseTime=true&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Failed connect DB:", err)
	}

	// 🔥 Ping untuk memastikan koneksi benar-benar hidup
	if err = db.Ping(); err != nil {
		log.Fatal("❌ DB not responding:", err)
	}

	// 🔥 OPTIONAL tapi bagus untuk production
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	DB = db

	log.Println("✅ Database connected")
}
