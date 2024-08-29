package handlers

import (
	"encoding/json"
	"net/http"

	"restapi/server/db"

	"golang.org/x/crypto/bcrypt"
)

// Benutzerstruktur für die Erstellung neuer Benutzer
type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // z.B. "admin" oder "user"
}

// Handler für die Erstellung eines neuen Benutzers
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Überprüfen, ob die Methode POST ist
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Methode erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// JSON-Daten aus der Anfrage lesen
	var newUser NewUser
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Ungültige JSON-Daten", http.StatusBadRequest)
		return
	}

	// Überprüfen der Rolle (sollte "admin" oder "user" sein)
	if newUser.Role != "admin" && newUser.Role != "user" {
		http.Error(w, "Ungültige Rolle", http.StatusBadRequest)
		return
	}

	// Passwort hashen
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Fehler beim Hashen des Passworts", http.StatusInternalServerError)
		return
	}

	// SQL-Abfrage zum Einfügen eines neuen Benutzers
	query := "INSERT INTO users (username, password, role) VALUES (?, ?, ?)"
	_, err = db.GetDB().Exec(query, newUser.Username, hashedPassword, newUser.Role)
	if err != nil {
		http.Error(w, "Fehler beim Einfügen des Benutzers in die Datenbank", http.StatusInternalServerError)
		return
	}

	// Erfolgreiche Antwort senden
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Benutzer erfolgreich erstellt"))
}

// Handler für die einmalige Erstellung eines Admin-Benutzers
func CreateFirstAdmin(w http.ResponseWriter, r *http.Request) {
	// Überprüfen, ob die Methode POST ist
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Methode erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Überprüfen, ob bereits ein Admin existiert
	var count int
	err := db.GetDB().QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		http.Error(w, "Fehler bei der Überprüfung vorhandener Admin-Benutzer", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Ein Admin-Benutzer existiert bereits", http.StatusForbidden)
		return
	}

	// JSON-Daten aus der Anfrage lesen
	var newUser NewUser
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Ungültige JSON-Daten", http.StatusBadRequest)
		return
	}

	// Überprüfen der Rolle (nur "admin" erlaubt)
	if newUser.Role != "admin" {
		http.Error(w, "Nur Admin-Rolle erlaubt", http.StatusBadRequest)
		return
	}

	// Passwort hashen
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Fehler beim Hashen des Passworts", http.StatusInternalServerError)
		return
	}

	// SQL-Abfrage zum Einfügen eines neuen Benutzers
	query := "INSERT INTO users (username, password, role) VALUES (?, ?, ?)"
	_, err = db.GetDB().Exec(query, newUser.Username, hashedPassword, newUser.Role)
	if err != nil {
		http.Error(w, "Fehler beim Einfügen des Benutzers in die Datenbank", http.StatusInternalServerError)
		return
	}

	// Erfolgreiche Antwort senden
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Admin-Benutzer erfolgreich erstellt"))
}
