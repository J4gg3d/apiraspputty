package models

// Order repräsentiert die Struktur der Daten aus der Tabelle dt_order
type Order struct {
	ID         int    `json:"id"`
	Kunde      int    `json:"kunde"`
	Kundenname string `json:"kundenname"`
}
