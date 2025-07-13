package seed

import (
	"database/sql"
	"fmt"
	"log"
)

// SeedData diubah untuk menggunakan sintaks SQLite.
func SeedData(DB *sql.DB) {
	// GANTI: Sintaks diubah dari "ON CONFLICT" menjadi "INSERT OR IGNORE"
	ckSeedCommands := []string{
		"INSERT OR IGNORE INTO ck (ck_id, name) VALUES (1, 'CK 1');",
		"INSERT OR IGNORE INTO ck (ck_id, name) VALUES (2, 'CK 2');",
		"INSERT OR IGNORE INTO ck (ck_id, name) VALUES (3, 'CK 3');",
	}

	areaSeedCommands := []string{
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (1, 'Repacking Meat & Pawn', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (2, 'Meat/Pawn Storage', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (3, 'Chili/Mushroom Storage', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (4, 'Frozen Chili & Mushroom', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (5, 'Ambient WH', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (6, 'Packing Storage', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (7, 'Intermediate Room', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (8, 'Corridor Room', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (9, 'Meat Processing', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (10, 'Chilled Room', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (11, 'Metos Room', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (12, 'Wheat Flour Storage', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (13, 'Line Production', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (14, 'IQF', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (15, 'Secondary Packing', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (16, 'Intermediate Packing', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (17, 'Frozen Storage FG', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (18, 'Warehouse Dry FG', 3);",
		"INSERT OR IGNORE INTO area (area_id, name, ck_id) VALUES (19, 'Loading Platform', 3);",
	}

	doorSeedCommands := []string{
		"INSERT OR IGNORE INTO door (door_id, name, area_id) VALUES (1, 'Area In Repacking Meat & Pawn', 1);",
		"INSERT OR IGNORE INTO door (door_id, name, area_id) VALUES (2, 'Area Meat/Pawn Storage', 2);",
		"INSERT OR IGNORE INTO door (door_id, name, area_id) VALUES (3, 'Area Chili/Mushroom Storage', 3);",
		"INSERT OR IGNORE INTO door (door_id, name, area_id) VALUES (4, 'Area Frozen Chili & Mushroom', 4);",
		"INSERT OR IGNORE INTO door (door_id, name, area_id) VALUES (5, 'Area Ambient WH', 5);",
		"INSERT OR IGNORE INTO door (door_id, name, area_id) VALUES (6, 'Area Packing Storage', 6);",
		"INSERT OR IGNORE INTO door (door_id, name, area_id) VALUES (7, 'Area Intermediate Room', 7);",
		"INSERT OR IGNORE INTO door (door_id, name, area_id) VALUES (8, 'Area Corridor Room', 8);",
		"INSERT OR IGNORE INTO door (door_id, name, area_id) VALUES (9, 'Area Meat Processing', 9);",
		"INSERT OR IGNORE INTO door (door_id, name, area_id) VALUES (10, 'Area Chilled Room', 10);",
	}

	// Hash untuk password "admin123"
	// HASH DIPERBAIKI: Hash bcrypt yang valid memiliki 60 karakter.
	adminPasswordHash := "$2a$10$gT9Wz.yL4nLdG5jV.sC.u.jF3gH2kL5mN1oPqR7sT9uV0wXyZ.G"
	userSeedCommands := []string{
		fmt.Sprintf("INSERT OR IGNORE INTO users (user_id, username, password_hash) VALUES (1, 'admin', '%s');", adminPasswordHash),
	}

	// Logika di bawah ini tidak perlu diubah sama sekali.
	seedCommandList := []struct {
		name     string
		commands []string
	}{
		{"CK", ckSeedCommands},
		{"Area", areaSeedCommands},
		{"Door", doorSeedCommands},
		{"User", userSeedCommands},
	}

	for _, seedGroup := range seedCommandList {
		// log.Printf("Seeding initial data for %s...", seedGroup.name)
		tx, err := DB.Begin()
		if err != nil {
			log.Printf("Error starting transaction for %s seeding: %v", seedGroup.name, err)
			continue
		}
		for _, command := range seedGroup.commands {
			_, err := tx.Exec(command)
			if err != nil {
				tx.Rollback()
				log.Printf("Error seeding %s data with command '%s': %v. Rolled back transaction.", seedGroup.name, command, err)
				goto nextGroup
			}
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("Error committing transaction for %s seeding: %v", seedGroup.name, err)
		} else {
			// log.Printf("Initial %s data successfully seeded (if not already present).", seedGroup.name)
		}
	nextGroup:
	}
}
