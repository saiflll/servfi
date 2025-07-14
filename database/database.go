package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" 
)

var DB *sql.DB


func InitDB() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("❌ DB_PATH tidak ditemukan di file .env")
	}

	var err error
	
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("❌ Gagal membuka database SQLite: %v", err)
	}

	
	if err = DB.Ping(); err != nil {
		log.Fatalf("❌ Gagal terhubung ke database: %v", err)
	}

	log.Println("🚀 Database SQLite berhasil terhubung.")
	
	
}


func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("🔌 Koneksi database ditutup.")
	}
}


func GetDB() *sql.DB {
	return DB
}
