// Package main är huvudpaketet för e-Arkive backend-servern
package main

import (
	"bytes"
	"context"
	"database/sql"
	"graphql-backend/graph"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

// =============================================
// ========== GLOBALA VARIABLER ===============
// =============================================

// Globala variabler
var db *sql.DB // Databasanslutningen som delas genom hela applikationen

// Standardport för servern om ingen annan specificerats
const defaultPort = "8080"

// =============================================
// ========== HJÄLPSTRUKTURER ================
// =============================================

// responseLogger används för att fånga upp och logga svaren från GraphQL-servern
type responseLogger struct {
	http.ResponseWriter
	body bytes.Buffer
}

func (rl *responseLogger) Write(b []byte) (int, error) {
	rl.body.Write(b)
	return rl.ResponseWriter.Write(b)
}

// =============================================
// ========== DATABAS-FUNKTIONER =============
// =============================================

// initDB initierar anslutningen till SQLite-databasen
// Skapar en ny databasfil om den inte redan finns
func initDB() {
	var err error
	connString := "./e-Arkive.db"
	db, err = sql.Open("sqlite3", connString)
	if err != nil {
		log.Fatalf("Failed to connect to SQLite database: %v", err)
	}

	// Verifiera databasanslutningen
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping SQLite database: %v", err)
	}

	// Skapa tabellerna om de inte redan finns
	createTables()

	log.Println("Connected to SQLite database successfully!")
}

// createTables skapar alla nödvändiga tabeller i databasen om de inte redan finns
func createTables() {
	// Läs SQL från filen
	sqlBytes, err := os.ReadFile("update_database.sql")
	if err != nil {
		log.Fatalf("Failed to read update_database.sql: %v", err)
	}

	// Exekvera SQL-skriptet
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("Failed to execute SQL script: %v", err)
	}

	log.Println("Database tables created or updated successfully")
}

// =============================================
// ========== LOGGNING OCH VERKTYG ===========
// =============================================

// logRequest loggar alla inkommande HTTP-förfrågningar
// Innehåller metod, sökväg och klientens IP-adress
func logRequest(r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
}

// logAction loggar viktiga händelser i systemet
func logAction(action string) {
	log.Printf("[ACTION] %s", action)
}

// getLocalIP hämtar serverns lokala IP-adress
// Används för att visa korrekt serveradress i loggarna
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Error getting local IP: %v", err)
		return "localhost"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}

	return "localhost"
}

// =============================================
// ========== ENDPOINT-HANTERARE =============
// =============================================

// setupGraphQLHandler konfigurerar och returnerar GraphQL-servern
func setupGraphQLHandler() *handler.Server {
	// Skapar en ny GraphQL-server med vår schema och resolver
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(db)}))

	// Konfigurerar tillåtna transportmetoder
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	// Konfigurerar cache för query-optimering
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	// Aktiverar GraphQL-tillägg
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

// setupQueryEndpoint konfigurerar /query-endpointen som hanterar GraphQL-förfrågningar
func setupQueryEndpoint(srv *handler.Server) {
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		logAction("GraphQL query received")

		// Extract JWT token from Authorization header
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Try Authenticate header if Authorization is not present
			authHeader = r.Header.Get("Authenticate")
		}

		if authHeader != "" {
			// Check if the header contains a Bearer token
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				// Add token to context using both keys for compatibility
				ctx = context.WithValue(ctx, "Authorization", token)
				ctx = context.WithValue(ctx, "Authenticate", token)
				// Create a new request with the updated context
				r = r.WithContext(ctx)
			}
		}

		responseRecorder := &responseLogger{ResponseWriter: w}
		srv.ServeHTTP(responseRecorder, r)
		log.Printf("Response sent: %s", responseRecorder.body.String())
	})
}

// setupStaticEndpoints konfigurerar ändpunkter för statiska resurser (GraphiQL, sandbox etc.)
func setupStaticEndpoints() {
	http.HandleFunc("/graphiql", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		http.ServeFile(w, r, "graphiql.html")
	})

	http.Handle("/sandbox", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(sandboxHTML)
	}))
}

// =============================================
// ========== HUVUDFUNKTION ==================
// =============================================

// main är huvudfunktionen som startar servern
func main() {
	// Initierar databasen
	initDB()

	// Konfigurerar serverporten
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Konfigurerar GraphQL-servern
	srv := setupGraphQLHandler()

	// Konfigurerar endpoints
	setupQueryEndpoint(srv)
	setupStaticEndpoints()

	localIP := getLocalIP()

	// Konfigurerar CORS för att tillåta anrop från frontend
	handlerWithCORS := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(http.DefaultServeMux)

	// Loggar serverinformation
	log.Printf("Server is running at http://%s:%s/query", localIP, port)
	log.Printf("GraphiQL is available at http://%s:%s/graphiql", localIP, port)
	log.Printf("Sandbox is available at http://%s:%s/sandbox", localIP, port)

	log.Printf("Server is starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithCORS))
}

// =============================================
// ========== STATISKA RESURSER ==============
// =============================================

// sandboxHTML innehåller HTML för GraphQL Playground
var sandboxHTML = []byte(`
<!DOCTYPE html>
<html lang="en">
<body style="margin: 0; overflow-x: hidden; overflow-y: hidden">
<div id="sandbox" style="height:100vh; width:100vw;"></div>
<script src="https://embeddable-sandbox.cdn.apollographql.com/_latest/embeddable-sandbox.umd.production.min.js"></script>
<script>
 new window.EmbeddedSandbox({
   target: "#sandbox",
   // Pass through your server href if you are embedding on an endpoint.
   // Otherwise, you can pass whatever endpoint you want Sandbox to start up with here.
   initialEndpoint: "http://localhost:8080/query",
 });
 // advanced options: https://www.apollographql.com/docs/studio/explorer/sandbox#embedding-sandbox
</script>
</body>
</html>`)
