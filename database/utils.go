package database

import (
	"log"
)

func (db *db) IsBootstrapped() bool {
	if db == nil || db.Tags == nil || db.Topics == nil {
		return false
	}
	return true
}

func (db *db) runMigrations() error {
	log.Println("🚸 Running migrations")
	//    _, err := db.Exec("CREATE TABLE IF NOT EXISTS tags (id INTEGER PRIMARY KEY, name TEXT)")
	//    if err != nil {
	//        return fmt.Errorf("database.runMigrations: %q\n", err)
	//    }
	//    _, err = db.Exec("CREATE TABLE IF NOT EXISTS topics (id INTEGER PRIMARY KEY, name TEXT)")
	//    if err != nil {
	//        return fmt.Errorf("database.runMigrations: %q\n", err)
	//    }
	return nil
}
