package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"restapi/server/db"
	"restapi/server/middleware"
)

// Order-Struktur
type Order struct {
	ID         int    `json:"id"` // ID der Order hinzufügen
	Kunde      int    `json:"kunde"`
	Kundenname string `json:"kundenname"`
	UserID     int    `json:"user_id"` // UserID hinzufügen
}

// @Summary Get all orders
// @Description Retrieve all orders for the authenticated user
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {array} Order "List of orders"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /orders [get]
func GetOrders(w http.ResponseWriter, r *http.Request) {
	// Benutzer-ID aus dem Kontext holen
	userID, ok := r.Context().Value(middleware.UserKey).(int)
	if !ok {
		http.Error(w, "Benutzer-ID nicht gefunden", http.StatusUnauthorized)
		return
	}

	// SQL-Abfrage, um alle Orders für diesen Benutzer abzurufen
	rows, err := db.GetDB().Query("SELECT ID, Kunde, Kundenname FROM dt_order WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Fehler bei der Abfrage der Datenbank", http.StatusInternalServerError)
		log.Printf("Fehler bei der Abfrage: %v", err)
		return
	}
	defer rows.Close()

	// Ergebnis in eine Slice von Order-Strukturen einlesen
	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.Kunde, &order.Kundenname)
		if err != nil {
			http.Error(w, "Fehler beim Lesen der Datenbankzeilen", http.StatusInternalServerError)
			log.Printf("Fehler beim Lesen der Zeile: %v", err)
			return
		}
		orders = append(orders, order)
	}

	// Überprüfen auf Fehler bei der Verarbeitung der Abfrageergebnisse
	if err = rows.Err(); err != nil {
		http.Error(w, "Fehler bei der Verarbeitung der Datenbankergebnisse", http.StatusInternalServerError)
		log.Printf("Fehler bei der Verarbeitung der Zeilen: %v", err)
		return
	}

	// JSON-Antwort erstellen und senden
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, "Fehler beim Erstellen der JSON-Antwort", http.StatusInternalServerError)
		log.Printf("Fehler beim JSON-Encoding: %v", err)
	}
}

// @Summary Create a new order
// @Description Create a new order in the system
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body Order true "Order to create"
// @Success 201 {string} string "Order successfully created"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /orders/create [post]
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Überprüfen, ob die Methode POST ist
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Methode erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Benutzer-ID aus dem Kontext holen
	userID, ok := r.Context().Value(middleware.UserKey).(int)
	if !ok {
		http.Error(w, "Benutzer-ID nicht gefunden", http.StatusUnauthorized)
		return
	}

	// JSON-Daten aus der Anfrage lesen
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Ungültige JSON-Daten", http.StatusBadRequest)
		log.Printf("Fehler beim Decodieren der JSON-Daten: %v", err)
		return
	}

	// SQL-Abfrage zum Einfügen einer neuen Order
	query := "INSERT INTO dt_order (Kunde, Kundenname, user_id) VALUES (?, ?, ?)"
	_, err = db.GetDB().Exec(query, order.Kunde, order.Kundenname, userID)
	if err != nil {
		http.Error(w, "Fehler beim Einfügen in die Datenbank", http.StatusInternalServerError)
		log.Printf("Fehler beim Einfügen: %v", err)
		return
	}

	// Erfolgreiche Antwort senden
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Order erfolgreich erstellt")
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// Überprüfen, ob die Methode DELETE ist
	if r.Method != http.MethodDelete {
		http.Error(w, "Nur DELETE-Methode erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Abrufen der ID aus den URL-Parametern
	queryParams := r.URL.Query()
	idParam := queryParams.Get("id")
	if idParam == "" {
		http.Error(w, "Order ID ist erforderlich", http.StatusBadRequest)
		return
	}

	// Konvertieren der ID in einen Integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Ungültige Order ID", http.StatusBadRequest)
		log.Printf("Ungültige Order ID: %v", err)
		return
	}

	// SQL-Abfrage zum Löschen einer Order
	query := "DELETE FROM dt_order WHERE ID = ?"
	result, err := db.GetDB().Exec(query, id)
	if err != nil {
		http.Error(w, "Fehler beim Löschen aus der Datenbank", http.StatusInternalServerError)
		log.Printf("Fehler beim Löschen: %v", err)
		return
	}

	// Überprüfen, ob eine Zeile gelöscht wurde
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Fehler beim Abrufen der betroffenen Zeilenanzahl", http.StatusInternalServerError)
		log.Printf("Fehler beim Abrufen der betroffenen Zeilenanzahl: %v", err)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Keine Order mit dieser ID gefunden", http.StatusNotFound)
		return
	}

	// Erfolgreiche Antwort senden
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order erfolgreich gelöscht")
}
