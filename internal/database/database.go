package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite" // GANTI: Menggunakan driver SQLite3

	"IoTT/internal/seed" // Asumsi path ini tetap sama
)

// Variabel global tetap sama
var DB *sql.DB
var (
	AreaNames     map[int]string
	DoorNames     map[int]string
	DoorToAreaMap map[int]int
)

// InitDB diubah untuk terhubung ke file SQLite dari .env
func InitDB() {
	// MENGGUNAKAN .env yang sudah di-load di main.go
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("❌ Error: DB_PATH tidak ditemukan di environment.")
	}

	var err error
	// GANTI: Membuka koneksi ke file SQLite
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("❌ Error saat membuka koneksi ke database SQLite: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("❌ Error saat ping database SQLite: %v", err)
	}

	log.Println("✅ Berhasil terhubung ke database SQLite.")

	// Menjalankan semua langkah inisialisasi
	createTables()
	log.Println("⚙️  Memulai proses seeding database...")
	seed.SeedData(DB)
	log.Println("✅ Seeding database selesai.")
	LoadLookupData()
}

// createTables disesuaikan untuk sintaks SQLite
func createTables() {
	// GANTI: Sintaks `SERIAL PRIMARY KEY` diubah menjadi `INTEGER PRIMARY KEY AUTOINCREMENT`
	commands := []string{
		`CREATE TABLE IF NOT EXISTS ck (
            ck_id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE
        );`,
		`CREATE TABLE IF NOT EXISTS area (
            area_id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            ck_id INTEGER NOT NULL,
            FOREIGN KEY (ck_id) REFERENCES ck (ck_id) ON DELETE CASCADE
        );`,
		`CREATE TABLE IF NOT EXISTS door (
            door_id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            area_id INTEGER NOT NULL,
            FOREIGN KEY (area_id) REFERENCES area (area_id) ON DELETE CASCADE
        );`,
		`CREATE TABLE IF NOT EXISTS temp (
            temp_id INTEGER PRIMARY KEY AUTOINCREMENT,
            value REAL NOT NULL,
            area_id INTEGER NOT NULL,
            no INTEGER NOT NULL,
            ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (area_id) REFERENCES area (area_id) ON DELETE CASCADE
        );`,
		`CREATE TABLE IF NOT EXISTS rh (
            rh_id INTEGER PRIMARY KEY AUTOINCREMENT,
            value REAL NOT NULL,
            area_id INTEGER NOT NULL,
            no INTEGER NOT NULL,
            ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (area_id) REFERENCES area (area_id) ON DELETE CASCADE
        );`,
		`CREATE TABLE IF NOT EXISTS prox (
            prox_id INTEGER PRIMARY KEY AUTOINCREMENT,
            value INTEGER NOT NULL,
            door_id INTEGER NOT NULL,
            ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (door_id) REFERENCES door (door_id) ON DELETE CASCADE
        );`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, command := range commands {
		_, err := DB.Exec(command)
		if err != nil {
			log.Fatalf("❌ Error membuat tabel dengan perintah %s: %v", command, err)
		}
	}
	log.Println("✅ Struktur tabel database berhasil diperiksa/dibuat.")
}

// LoadLookupData tidak perlu diubah, karena menggunakan query SQL standar
func LoadLookupData() {
	AreaNames = make(map[int]string)
	DoorNames = make(map[int]string)
	DoorToAreaMap = make(map[int]int)

	rows, err := DB.Query("SELECT area_id, name FROM area")
	if err != nil {
		log.Printf("Error memuat nama area: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Printf("Error scanning baris area: %v", err)
			continue
		}
		AreaNames[id] = name
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error setelah iterasi baris area: %v", err)
	}

	doorRows, err := DB.Query("SELECT door_id, name, area_id FROM door")
	if err != nil {
		log.Printf("Error memuat nama pintu: %v", err)
		return
	}
	defer doorRows.Close()

	for doorRows.Next() {
		var id, areaID int
		var name string
		if err := doorRows.Scan(&id, &name, &areaID); err != nil {
			log.Printf("Error scanning baris pintu: %v", err)
			continue
		}
		DoorNames[id] = name
		DoorToAreaMap[id] = areaID
	}
	if err = doorRows.Err(); err != nil {
		log.Printf("Error setelah iterasi baris pintu: %v", err)
	}
	log.Printf("✔️ Berhasil memuat %d nama area dan %d nama pintu ke lookup map.", len(AreaNames), len(DoorNames))
}

// Helper functions tidak perlu diubah
func GetDB() *sql.DB {
	return DB
}

func GetAreaName(areaID int) string {
	if name, ok := AreaNames[areaID]; ok {
		return name
	}
	return fmt.Sprintf("Area ID %d (Nama tidak ditemukan)", areaID)
}

func GetDoorInfo(doorID int) (doorName string, areaID int, areaName string) {
	var ok bool
	doorName, ok = DoorNames[doorID]
	if !ok {
		doorName = fmt.Sprintf("Pintu ID %d (Nama tidak ditemukan)", doorID)
	}

	areaID, ok = DoorToAreaMap[doorID]
	if !ok {
		areaName = "Area tidak diketahui untuk pintu ini"
		return doorName, 0, areaName
	}

	areaName = GetAreaName(areaID)
	return doorName, areaID, areaName
}

// CloseDB adalah fungsi baru untuk menutup koneksi database dengan rapi.
// Panggil dengan `defer database.CloseDB()` di `main.go`.
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
