package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Schlüssel für Benutzerkontext
type key int

const (
	UserKey key = iota
)

var jwtKey = []byte("geheimesSchluessel")

// Middleware zur Überprüfung des JWT-Tokens und der Berechtigungen
func AuthMiddleware(allowedRoles []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Token aus dem Authorization-Header abrufen
		tokenString := r.Header.Get("Authorization")

		// Prüfen, ob der Header vorhanden ist und das Token das richtige Format hat
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			http.Error(w, "Token erforderlich", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Token validieren
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Ungültiges Token", http.StatusUnauthorized)
			return
		}

		// Überprüfen, ob die Rolle erlaubt ist
		userRole, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Rolle nicht gefunden", http.StatusForbidden)
			return
		}

		roleAllowed := false
		for _, role := range allowedRoles {
			if role == userRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			http.Error(w, "Keine Berechtigung für diese Operation", http.StatusForbidden)
			return
		}

		// Benutzer-ID als Integer aus den Claims extrahieren
		userID, ok := claims["user_id"].(float64) // JWT speichert numerische Werte als float64
		if !ok {
			http.Error(w, "Benutzer-ID nicht gefunden", http.StatusForbidden)
			return
		}

		// Benutzer in den Kontext setzen (hier als int konvertiert)
		ctx := context.WithValue(r.Context(), UserKey, int(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
