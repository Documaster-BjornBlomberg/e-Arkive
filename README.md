# e-Arkive

e-Arkive är en arkivapplikation som består av en backend och en frontend. Applikationen låter användare ladda upp dokument och spara dem i en databas tillsammans med relevant metadata.

## Backend

### Tekniker
- **Språk:** Go  
- **Databas:** SQLite (lagras som en fil)  
- **API:** GraphQL (med hjälp av `gqlgen`)
- **Autentisering:** JWT (JSON Web Tokens)

### Funktioner
- Hanterar filuppladdningar och metadata.  
- Använder SQLite som en lättviktig databas.  
- Tillhandahåller ett GraphQL API för frontend-kommunikation.
- Stödjer användarautentisering och behörighetshantering.

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

### Autentisering och användarhantering

e-Arkive använder ett JWT-baserat (JSON Web Token) autentiseringssystem med följande funktioner:

#### Användarkonton
- **Standardanvändare:** Systemet levereras med en fördefinierad admin-användare (användarnamn: `admin`, lösenord: `admin`)
- **Registrering:** Nya användare kan registrera sig via API:et med ett användarnamn och lösenord
- **Lösenordshantering:** Lösenord lagras säkert med bcrypt-hashing
- **Användargrupper:** Användare kan tillhöra grupper för behörighetshantering

#### Autentiseringsflöde
1. **Inloggning:** Användaren skickar användarnamn och lösenord till `/query`-endpunkten
2. **Token-generering:** Vid framgångsrik inloggning genererar servern en JWT som är giltig i 7 dagar
3. **Auktorisering:** Efterföljande API-anrop inkluderar token i Authorization-headern (`Bearer <token>`)
4. **Utloggning:** Tokens kan ogiltigförklaras vid utloggning genom att lägga till dem i en svartlista

#### Behörighetsmodell
e-Arkive använder en nodbaserad hierarkisk behörighetsmodell:

- **Binära behörigheter:** Rättigheter representeras som binära flaggor
   - READ (1): Rätt att läsa/visa en nod och dess innehåll
   - MODIFY (2): Rätt att ändra en nod eller dess innehåll
   - DELETE (4): Rätt att ta bort en nod eller dess innehåll

- **Behörighetsnivåer:**
   - **Nodens ägare:** Har alltid fulla rättigheter (7 = läs + ändra + ta bort)
   - **Gruppmedlemmar:** Om noden ägs av en grupp, ärver gruppens medlemmar nodens behörigheter
   - **Övriga användare:** Har de behörigheter som är definierade i nodens permissions-fält

### Noder och filsystem

e-Arkive organiserar innehåll i en hierarkisk nodstruktur:

- **Root-nod:** Systemet har en rotnod som fungerar som utgångspunkt i hierarkin
- **Noder:** En nod kan representera en mapp i filsystemet
- **Filer:** Filer kan kopplas till noder och innehålla metadata
- **Behörigheter:** Varje nod har specifika behörighetsinställningar

### Testa API:et

För att testa GraphQL API:et kan du använda det inbyggda webbgränssnittet:

1. Öppna `http://localhost:8080/sandbox` i en webbläsare
2. Utför inloggning för att få en token:
   ```graphql
   mutation {
     login(username: "admin", password: "admin") {
       token
       user {
         id
         username
       }
     }
   }
   ```
3. Kopiera den returnerade token och klicka på "Headers" längst ned:
   ```json
   {
     "Authorization": "Bearer ditt_token_här"
   }
   ```
4. Nu kan du utföra auktoriserade anrop, till exempel:
   ```graphql
   query {
     me {
       id
       username
     }
   }
   ```

#### Exempel på API-anrop

**Hämta användarens inställningar:**
```graphql
query {
  getUserSettings {
    key
    value
  }
}
```

**Hämta filhierarkin:**
```graphql
query {
  getRootNodes {
    id
    name
    children {
      id
      name
    }
  }
}
```

**Ladda upp en fil:**
```graphql
mutation {
  saveFile(input: {
    name: "exempel.txt",
    size: 123,
    contentType: "text/plain",
    fileData: "base64kodad_fildata",
    metadata: [
      { key: "author", value: "Anders Andersson" },
      { key: "description", value: "Exempelfil" }
    ],
    nodeId: "1"
  }) {
    id
    name
  }
}
```

### Databasstruktur

e-Arkive använder SQLite för att lagra alla data. Huvudtabellerna är:

- **users:** Användarinformation och inloggningsuppgifter
- **user_settings:** Användarspecifika inställningar (t.ex. tema)
- **groups:** Användargrupper för behörighetshantering
- **group_members:** Kopplingar mellan användare och grupper
- **nodes:** Hierarkisk struktur som representerar mappträdet
- **files:** Filinformation och binärdata
- **metadata:** Metadata kopplad till filer som nyckel-värde-par

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
│   ├── auth/              # Autentisering och behörigheter
│   │   ├── auth.go       # JWT-hantering, token-validering
│   │   └── permissions.go # Behörighetshantering och kontroller
│   ├── datastore/        # Databasoperationer och dataåtkomst
│   │   ├── db.go         # Databasanslutning
│   │   ├── file_store.go # Filhanteringsoperationer
│   │   ├── node_store.go # Nodhanteringsoperationer
│   │   └── user_store.go # Användarhanteringsoperationer
│   ├── graph/             # GraphQL-schema och resolvers
│   │   ├── model/        # GraphQL-datamodeller
│   │   ├── file_resolver.go # Filtypsresolver
│   │   ├── node_resolver.go # Nodtypsresolver
│   │   ├── user_resolver.go # Användartypsresolver
│   │   ├── resolver.go   # Huvudresolver med beroendeinjektion
│   │   ├── schema.resolvers.go # Implementationer av queries och mutations
│   │   └── schema.graphqls # GraphQL-schema
│   ├── util/              # Hjälpfunktioner
│   │   ├── cycle_detection.go # Upptäcker cykler i nodhierarkier
│   │   └── log.go        # Loggningsfunktioner
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
