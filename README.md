# e-Arkive

e-Arkive är en arkivapplikation som består av en backend och en frontend. Applikationen låter användare ladda upp dokument och spara dem i en databas tillsammans med relevant metadata.

## Backend

### Tekniker
- **Språk:** Go
- **Databas:** SQLite (lagras som en fil)
- **API:** GraphQL (med hjälp av `gqlgen`)

### Funktioner
- Hanterar filuppladdningar och metadata.
- Använder SQLite som en lättviktig databas.
- Tillhandahåller ett GraphQL API för frontend-kommunikation.

### Hur man startar backend
1. Se till att du har Go installerat.
2. Installera nödvändiga beroenden:
   ```bash
   go mod tidy
   ```
3. Starta backend:
   ```bash
   go run .
   ```
4. Backend körs på `http://localhost:8080`. Du kan använda GraphQL Playground på `http://localhost:8080/`.

## Frontend

### Tekniker
- **Ramverk:** Vue 3
- **Byggverktyg:** Vite

### Funktioner
- Låter användare ladda upp dokument via ett enkelt gränssnitt.
- Kommunicerar med backend via GraphQL API.

### Hur man startar frontend
1. Se till att du har Node.js och npm installerat.
2. Installera beroenden:
   ```bash
   npm install
   ```
3. Starta frontend:
   ```bash
   npm run dev
   ```
4. Frontend körs på `http://localhost:5173`.

## Krav

- **Backend**:
  - Go 1.20 eller senare.
  - SQLite (stöds via `go-sqlite3`).
  - TDM-GCC eller annan C-kompilator (för att bygga `go-sqlite3`).

- **Frontend**:
  - Node.js 16 eller senare.
  - npm 7 eller senare.

## Projektstruktur

```plaintext
...existing code...
├── backend/               # Backend-kod
│   ├── graph/             # GraphQL-schema och resolvers
│   ├── e-Arkive.db        # SQLite-databasfil
│   └── server.go          # Huvudfil för backend
├── frontend/              # Frontend-kod
│   ├── src/               # Vue-komponenter
│   └── vite.config.js     # Vite-konfiguration
└── README.md              # Projektbeskrivning
```
