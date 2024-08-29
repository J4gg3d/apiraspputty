package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Globale Datenbankvariable
var db *sql.DB

// Funktion zum Herstellen der Datenbankverbindung
func ConnectToDB() {
	var err error
	username := "user1"
	password := "harry_hirsch"
	hostname := "192.168.2.172"
	port := 3306
	dbname := "mariadb"

	// Verbindung zur Datenbank herstellen
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, hostname, port, dbname)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Fehler bei der Verbindung zur Datenbank: %v", err)
	}

	// Testen der Datenbankverbindung
	err = db.Ping()
	if err != nil {
		log.Fatalf("Kann keine Verbindung zur Datenbank herstellen: %v", err)
	}

	fmt.Println("Erfolgreich mit der Datenbank verbunden!")
}

// Funktion zum Schließen der Datenbankverbindung
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// Funktion zum Abrufen der DB-Instanz
func GetDB() *sql.DB {
	return db
}
