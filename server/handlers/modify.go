package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"restapi/server/db"
	"restapi/server/models"
)

// Handler für PUT-Anfragen zum Ändern einer Order
func ModifyOrder(w http.ResponseWriter, r *http.Request) {
	// Überprüfen, ob die Methode PUT ist
	if r.Method != http.MethodPut {
		http.Error(w, "Nur PUT-Methode erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// JSON-Daten aus der Anfrage lesen
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Ungültige JSON-Daten", http.StatusBadRequest)
		log.Printf("Fehler beim Decodieren der JSON-Daten: %v", err)
		return
	}

	// Überprüfen, ob die ID für das Update vorhanden ist
	if order.ID == 0 {
		http.Error(w, "Order ID ist erforderlich", http.StatusBadRequest)
		return
	}

	// SQL-Abfrage zum Aktualisieren einer Order
	query := "UPDATE dt_order SET Kunde = ?, Kundenname = ? WHERE ID = ?"
	result, err := db.GetDB().Exec(query, order.Kunde, order.Kundenname, order.ID)
	if err != nil {
		http.Error(w, "Fehler beim Aktualisieren der Datenbank", http.StatusInternalServerError)
		log.Printf("Fehler beim Aktualisieren: %v", err)
		return
	}

	// Überprüfen, ob eine Zeile aktualisiert wurde
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
	fmt.Fprintf(w, "Order erfolgreich aktualisiert")
}
