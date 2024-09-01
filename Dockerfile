# Verwenden Sie ein offizielles Go-Bild als Build-Image
FROM golang:1.20 AS builder

# Arbeitsverzeichnis festlegen
WORKDIR /app

# Kopieren der go.mod und go.sum aus dem Hauptverzeichnis
COPY go.mod go.sum ./

# Herunterladen der Abhängigkeiten
RUN go mod download

# Arbeitsverzeichnis auf das 'server'-Verzeichnis ändern
WORKDIR /app/server

# Kopieren des gesamten 'server'-Unterverzeichnisses in das Arbeitsverzeichnis
COPY server/ .

# Die Anwendung für Linux/ARM64 bauen
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o app .

# Verwenden Sie ein schlankes Image für die endgültige Ausführung
FROM alpine:latest  

WORKDIR /root/

# Kopieren Sie die gebaute Anwendung aus dem Build-Image
COPY --from=builder /app/server/app .

# Port 8080 für die Anwendung freigeben
EXPOSE 8080

# Starten der Anwendung
CMD ["./app"]
