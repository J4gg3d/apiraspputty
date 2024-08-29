package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"restapi/server/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Benutzerstruktur für einfache Authentifizierung
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Secret key (in a real application, keep this secret!)
var jwtKey = []byte("geheimesSchluessel")

// Handler für die Login-Route
func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ungültige Anfragedaten", http.StatusBadRequest)
		return
	}

	// Benutzer aus der Datenbank abrufen
	query := "SELECT id, username, password, role FROM users WHERE username = ?"
	row := db.GetDB().QueryRow(query, user.Username)

	var dbUser User
	err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password, &dbUser.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Ungültige Anmeldeinformationen", http.StatusUnauthorized)
		} else {
			http.Error(w, "Fehler beim Abrufen der Benutzerdaten", http.StatusInternalServerError)
		}
		return
	}

	// Passwort überprüfen
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Ungültige Anmeldeinformationen", http.StatusUnauthorized)
		return
	}

	// JWT erstellen mit Benutzer-ID und Rolle als Claim
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := jwt.MapClaims{
		"user_id":  dbUser.ID, // Benutzer-ID als Integer speichern
		"username": dbUser.Username,
		"role":     dbUser.Role,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Fehler beim Erstellen des Tokens", http.StatusInternalServerError)
		return
	}

	// Token als JSON zurückgeben
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
