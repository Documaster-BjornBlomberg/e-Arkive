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
1. Öppna en terminal i projektmappen `~\e-Arkive\graphql-backend`.
2. Installera nödvändiga beroenden:  
   ```bash
   go mod tidy
   ```
3. Starta backend:
   ```bash
   go run .
   ```
4. Backend körs på `http://localhost:8080`. Du kan använda GraphQL Playground på `http://localhost:8080/sandbox`.

## Frontend

### Tekniker
- **Ramverk:** Vue 3
- **Byggverktyg:** Vite
- **Arkitektur:** Atomic Design
- **Stilning:** CSS med variabelbaserade teman för ljust/mörkt läge

### Atomic Design
Projektet följer Atomic Design-principerna, vilket innebär att UI-komponenter är organiserade i ett hierarkiskt system:

- **Atomer:** Grundläggande byggstenar som knappar, ikoner och inmatningsfält
- **Molekyler:** Grupper av atomer som fungerar tillsammans (t.ex. sökfält, metadata-fält)
- **Organismer:** Mer komplexa komponenter som kombinerar flera molekyler (t.ex. fillistor, sidopaneler)
- **Mallar:** Sidlayouter utan faktiskt innehåll
- **Sidor:** Kompletta vyer med faktiskt innehåll

**Viktigt:** Vid implementering av nya funktioner, använd alltid befintliga komponenter från lägre nivåer innan nya skapas. Detta främjar återanvändning och konsekvent användargränssnitt.

### Huvudfunktioner
- Filuppladdning med metadata
- Fillistning och sökning
- Filvisning och nedladdning
- Metadata-hantering (lägga till, redigera, ta bort)
- Responsiv design som fungerar på olika skärmstorlekar
- Stöd för både ljust och mörkt tema

### Hur man startar frontend
1. Se till att du har Node.js och npm installerat.
2. Installera beroenden:
   ```powershell
   npm install
   ```
3. Starta frontend:
   ```powershell
   npm run dev
   ```
4. Frontend körs på `http://localhost:5173`.

### Utveckling på Windows 11
Projektet är konfigurerat för utveckling på Windows 11 med PowerShell som standardterminal:

1. Använd Visual Studio Code för bästa utvecklingsupplevelse
2. PowerShell-kommandon för vanliga uppgifter:
   ```powershell
   # Starta både backend och frontend
   ./start-projects.ps1
   
   # Installera nya npm-paket
   npm install paketnamn --save
   
   # Bygg för produktion
   npm run build
   ```

## Krav

- **Backend**:
  - Go 1.20 eller senare.
  - SQLite (stöds via `go-sqlite3`).
  - TDM-GCC eller annan C-kompilator (för att bygga `go-sqlite3`).

- **Frontend**:
  - Node.js 16 eller senare.
  - npm 7 eller senare.
  - Windows 11 med PowerShell

## Projektstruktur

```plaintext
├── graphql-backend/        # Backend-kod
│   ├── graph/             # GraphQL-schema och resolvers
│   │   ├── model/        # GraphQL-datamodeller
│   │   └── schema.graphqls # GraphQL-schema
│   ├── e-Arkive.db        # SQLite-databasfil
│   └── server.go          # Huvudfil för backend
├── src/                   # Frontend-kod
│   ├── components/        # Vue-komponenter enligt Atomic Design
│   │   ├── atoms/        # Grundläggande komponenter
│   │   ├── molecules/    # Sammansatta komponenter
│   │   ├── organisms/    # Komplexa komponenter
│   │   ├── templates/    # Sidlayouter
│   │   └── pages/        # Kompletta sidor
│   ├── composables/       # Vue 3 kompositionsfunktioner
│   ├── assets/            # Statiska tillgångar
│   ├── App.vue            # Huvudapplikationskomponent
│   └── main.js            # Applikationsstartpunkt
├── public/                # Statiska filer som kopieras direkt
├── vite.config.js         # Vite-konfiguration
└── README.md              # Projektbeskrivning (denna fil)
```

### Uppdatera GraphQL-koden

Om du gör ändringar i GraphQL-schemat eller resolvrarna, måste du generera om koden med `gqlgen`. Följ dessa steg:

1. Installera `gqlgen` och dess beroenden:
   ```bash
   go get github.com/99designs/gqlgen@v0.17.69
   ```

2. Generera om koden:
   ```bash
   go run github.com/99designs/gqlgen
   ```

3. Efter att koden har genererats om, kontrollera `schema.resolvers.go` för att säkerställa att alla resolvrar är korrekt implementerade. Om några resolvrar saknas eller behöver uppdateras, implementera dem manuellt.
