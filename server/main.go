package main

import (
	"fmt"
	"log"
	"net/http"

	// Importiert die generierte Swagger-Dokumentation
	"restapi/server/db"
	"restapi/server/handlers"
	"restapi/server/middleware"

	_ "restapi/server/docs" // Importieren der generierten Swagger-Dokumentation

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Verbindung zur Datenbank herstellen
	db.ConnectToDB()
	defer db.CloseDB()

	// Routen definieren
	http.HandleFunc("/api/login", handlers.Login) // Login-Route

	// Einmaliger Endpunkt zur Erstellung des ersten Admin-Benutzers
	http.HandleFunc("/api/users/create-first-admin", handlers.CreateFirstAdmin) // POST für ersten Admin

	// Geschützte Routen
	http.Handle("/api/orders", middleware.AuthMiddleware([]string{"user", "admin"}, http.HandlerFunc(handlers.GetOrders)))          // GET
	http.Handle("/api/orders/create", middleware.AuthMiddleware([]string{"user", "admin"}, http.HandlerFunc(handlers.CreateOrder))) // POST
	http.Handle("/api/orders/modify", middleware.AuthMiddleware([]string{"admin"}, http.HandlerFunc(handlers.ModifyOrder)))         // PUT
	http.Handle("/api/orders/delete", middleware.AuthMiddleware([]string{"admin"}, http.HandlerFunc(handlers.DeleteOrder)))         // DELETE
	http.Handle("/api/users/create", middleware.AuthMiddleware([]string{"admin"}, http.HandlerFunc(handlers.CreateUser)))           // POST für Benutzererstellung

	// Swagger UI Setup
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Server starten
	fmt.Println("Server läuft auf Port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
