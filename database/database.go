package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import driver SQLite3
)

var DB *sql.DB

// InitDB menginisialisasi koneksi ke database SQLite berdasarkan path dari .env.
func InitDB() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("❌ DB_PATH tidak ditemukan di file .env")
	}

	var err error
	// Membuka (atau membuat jika belum ada) file database SQLite
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("❌ Gagal membuka database SQLite: %v", err)
	}

	// Memeriksa apakah koneksi berhasil
	if err = DB.Ping(); err != nil {
		log.Fatalf("❌ Gagal terhubung ke database: %v", err)
	}

	log.Println("🚀 Database SQLite berhasil terhubung.")
	// Di sini Anda bisa menambahkan fungsi untuk membuat tabel jika belum ada (opsional)
	// createTables()
}

// CloseDB menutup koneksi database.
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("🔌 Koneksi database ditutup.")
	}
}

// GetDB adalah helper untuk mendapatkan instance DB yang aktif.
func GetDB() *sql.DB {
	return DB
}
