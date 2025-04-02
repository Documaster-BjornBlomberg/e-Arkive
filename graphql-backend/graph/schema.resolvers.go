package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.69

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"graphql-backend/graph/model"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SaveFile är resolvern för saveFile-fältet
// Hanterar uppladdning av nya filer och deras metadata till databasen
func (r *mutationResolver) SaveFile(ctx context.Context, input model.FileInput) (*model.File, error) {
	logAction("Received request to save file")
	log.Printf("Saving file: %s", input.Name)

	// Verifierar att databasanslutningen är aktiv
	if r.DB == nil {
		log.Printf("Database connection is nil")
		return nil, fmt.Errorf("internal server error: database connection is not initialized")
	}

	// Dekodar base64-data till binär form
	fileData, err := base64.StdEncoding.DecodeString(input.FileData)
	if err != nil {
		log.Printf("Error decoding file data: %v", err)
		return nil, fmt.Errorf("invalid file data: %v", err)
	}

	// Sparar filinformation och binär data i databasen
	result, err := r.DB.Exec(
		"INSERT INTO files (name, size, content_type, created_at, file_data) VALUES (?, ?, ?, datetime('now'), ?)",
		input.Name, input.Size, input.ContentType, fileData,
	)
	if err != nil {
		log.Printf("Error saving file to database: %v", err)
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	fileID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v", err)
		return nil, fmt.Errorf("failed to retrieve file ID: %v", err)
	}

	log.Printf("File saved to database with ID: %d", fileID)

	// Sparar metadata för filen
	for _, meta := range input.Metadata {
		_, err := r.DB.Exec(
			"INSERT INTO metadata (file_id, key, value) VALUES (?, ?, ?)",
			fileID, meta.Key, meta.Value,
		)
		if err != nil {
			log.Printf("Error saving metadata to database: %v", err)
			return nil, fmt.Errorf("failed to save metadata: %v", err)
		}
	}

	log.Printf("File and metadata saved successfully with ID: %d", fileID)

	// Konverterar metadata till rätt format för responsen
	metadata := make([]*model.Metadata, len(input.Metadata))
	for i, meta := range input.Metadata {
		metadata[i] = &model.Metadata{
			Key:   meta.Key,
			Value: meta.Value,
		}
	}

	return &model.File{
		ID:          fmt.Sprintf("%d", fileID),
		Name:        input.Name,
		Size:        input.Size,
		ContentType: input.ContentType,
		CreatedAt:   time.Now().Format(time.RFC3339),
		FileData:    &input.FileData, // Skickar tillbaka base64-kodad data
		Metadata:    metadata,
	}, nil
}

// deleteFile är resolvern för deleteFile-mutation
func (r *mutationResolver) DeleteFile(ctx context.Context, id string) (bool, error) {
	logAction(fmt.Sprintf("Attempting to delete file with ID: %s", id))

	if r.DB == nil {
		log.Printf("Database connection is nil")
		return false, fmt.Errorf("internal server error: database connection is not initialized")
	}

	// Ta bort filen och dess metadata från databasen
	result, err := r.DB.Exec("DELETE FROM files WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting file with ID %s: %v", id, err)
		return false, fmt.Errorf("failed to delete file: %v", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching affected rows for file ID %s: %v", id, err)
		return false, fmt.Errorf("failed to verify deletion: %v", err)
	}

	if affectedRows == 0 {
		log.Printf("No file found with ID: %s", id)
		return false, fmt.Errorf("file not found")
	}

	log.Printf("Successfully deleted file with ID: %s", id)
	return true, nil
}

// UpdateMetadata är resolvern för updateMetadata-mutation
// Uppdaterar metadata för en fil
func (r *mutationResolver) UpdateMetadata(ctx context.Context, fileID string, metadataInput []*model.MetadataInput) (*model.File, error) {
	logAction(fmt.Sprintf("Attempting to update metadata for file with ID: %s", fileID))

	if r.DB == nil {
		log.Printf("Database connection is nil")
		return nil, fmt.Errorf("internal server error: database connection is not initialized")
	}

	// Kontrollera att filen finns
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM files WHERE id = ?)", fileID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if file exists: %v", err)
		return nil, fmt.Errorf("failed to check if file exists: %v", err)
	}

	if !exists {
		log.Printf("No file found with ID: %s", fileID)
		return nil, fmt.Errorf("file not found")
	}

	// Starta en transaktion för att säkerställa att alla operationer lyckas eller misslyckas tillsammans
	tx, err := r.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Ta bort befintlig metadata för filen
	_, err = tx.Exec("DELETE FROM metadata WHERE file_id = ?", fileID)
	if err != nil {
		log.Printf("Error deleting existing metadata: %v", err)
		return nil, fmt.Errorf("failed to delete existing metadata: %v", err)
	}

	// Lägg till ny metadata
	for _, meta := range metadataInput {
		_, err := tx.Exec(
			"INSERT INTO metadata (file_id, key, value) VALUES (?, ?, ?)",
			fileID, meta.Key, meta.Value,
		)
		if err != nil {
			log.Printf("Error inserting metadata: %v", err)
			return nil, fmt.Errorf("failed to insert metadata: %v", err)
		}
	}

	// Commit transaktionen
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully updated metadata for file with ID: %s", fileID)

	// Hämta den uppdaterade filen för att returnera
	return (&queryResolver{r.Resolver}).GetFile(ctx, fileID)
}

// DeleteMetadata är resolvern för deleteMetadata-mutation
// Tar bort specifik metadata från en fil baserat på nycklar
func (r *mutationResolver) DeleteMetadata(ctx context.Context, fileID string, keys []string) (*model.File, error) {
	logAction(fmt.Sprintf("Attempting to delete metadata for file with ID: %s", fileID))

	if r.DB == nil {
		log.Printf("Database connection is nil")
		return nil, fmt.Errorf("internal server error: database connection is not initialized")
	}

	// Kontrollera att filen finns
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM files WHERE id = ?)", fileID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if file exists: %v", err)
		return nil, fmt.Errorf("failed to check if file exists: %v", err)
	}

	if !exists {
		log.Printf("No file found with ID: %s", fileID)
		return nil, fmt.Errorf("file not found")
	}

	// Starta en transaktion
	tx, err := r.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Ta bort metadata baserat på nycklar
	for _, key := range keys {
		_, err := tx.Exec("DELETE FROM metadata WHERE file_id = ? AND key = ?", fileID, key)
		if err != nil {
			log.Printf("Error deleting metadata with key %s: %v", key, err)
			return nil, fmt.Errorf("failed to delete metadata: %v", err)
		}
	}

	// Commit transaktionen
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully deleted metadata for file with ID: %s", fileID)

	// Hämta den uppdaterade filen för att returnera
	return (&queryResolver{r.Resolver}).GetFile(ctx, fileID)
}

// GetFiles är resolvern för getFiles-fältet
// Hämtar alla filer från databasen med tillhörande metadata
func (r *queryResolver) GetFiles(ctx context.Context) ([]*model.File, error) {
	logAction("Fetching all files from the database")

	// Hämtar alla filer från databasen
	rows, err := r.DB.Query("SELECT id, name, size, content_type, created_at, file_data FROM files")
	if err != nil {
		log.Printf("Error fetching files from database: %v", err)
		return nil, fmt.Errorf("failed to fetch files: %v", err)
	}
	defer rows.Close()

	var files []*model.File
	for rows.Next() {
		var file model.File
		var createdAt string
		var fileData []byte
		if err := rows.Scan(&file.ID, &file.Name, &file.Size, &file.ContentType, &createdAt, &fileData); err != nil {
			log.Printf("Error scanning file row: %v", err)
			return nil, fmt.Errorf("failed to scan file row: %v", err)
		}
		file.CreatedAt = createdAt
		if fileData != nil {
			encodedFileData := base64.StdEncoding.EncodeToString(fileData)
			file.FileData = &encodedFileData
		}

		// Hämtar metadata för varje fil
		metaRows, err := r.DB.Query("SELECT key, value FROM metadata WHERE file_id = ?", file.ID)
		if err != nil {
			log.Printf("Error fetching metadata for file ID %s: %v", file.ID, err)
			return nil, fmt.Errorf("failed to fetch metadata: %v", err)
		}
		defer metaRows.Close()

		var metadata []*model.Metadata
		for metaRows.Next() {
			var meta model.Metadata
			if err := metaRows.Scan(&meta.Key, &meta.Value); err != nil {
				log.Printf("Error scanning metadata row: %v", err)
				return nil, fmt.Errorf("failed to scan metadata row: %v", err)
			}
			metadata = append(metadata, &meta)
		}
		file.Metadata = metadata
		files = append(files, &file)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over file rows: %v", err)
		return nil, fmt.Errorf("failed to iterate over file rows: %v", err)
	}

	log.Printf("Successfully fetched %d files", len(files))
	return files, nil
}

// GetFile är resolvern för getFile-fältet
// Hämtar en specifik fil baserat på ID utan binärdata
func (r *queryResolver) GetFile(ctx context.Context, id string) (*model.File, error) {
	logAction(fmt.Sprintf("Fetching file with ID: %s", id))

	if r.DB == nil {
		err := fmt.Errorf("database connection is not initialized")
		log.Printf("GetFile failed: %v", err)
		return nil, err
	}

	// Hämtar filinformation från databasen utan binärdata
	var file model.File
	var createdAt string
	err := r.DB.QueryRow(`
		SELECT id, name, size, content_type, created_at
		FROM files WHERE id = ?`, id).Scan(
		&file.ID, &file.Name, &file.Size, &file.ContentType, &createdAt)

	if err == sql.ErrNoRows {
		errMsg := fmt.Sprintf("file not found with ID: %s", id)
		log.Printf("GetFile failed: %s", errMsg)
		return nil, fmt.Errorf(errMsg)
	} else if err != nil {
		log.Printf("GetFile failed: Error fetching file with ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to fetch file: %v", err)
	}

	file.CreatedAt = createdAt

	// Hämtar metadata för filen
	metaRows, err := r.DB.Query("SELECT key, value FROM metadata WHERE file_id = ?", file.ID)
	if err != nil {
		log.Printf("Error fetching metadata for file ID %s: %v", file.ID, err)
		return nil, fmt.Errorf("failed to fetch metadata: %v", err)
	}
	defer metaRows.Close()

	var metadata []*model.Metadata
	for metaRows.Next() {
		var meta model.Metadata
		if err := metaRows.Scan(&meta.Key, &meta.Value); err != nil {
			log.Printf("Error scanning metadata row: %v", err)
			return nil, fmt.Errorf("failed to scan metadata row: %v", err)
		}
		metadata = append(metadata, &meta)
	}
	file.Metadata = metadata

	log.Printf("Successfully fetched file %s (ID: %s)", file.Name, file.ID)
	return &file, nil
}

// DownloadFile är resolvern för downloadFile-fältet
// Hämtar en fil för nedladdning baserat på ID
func (r *queryResolver) DownloadFile(ctx context.Context, id string) (*model.File, error) {
	logAction(fmt.Sprintf("Attempting to download file with ID: %s", id))

	if r.DB == nil {
		err := fmt.Errorf("database connection is not initialized")
		log.Printf("Download failed: %v", err)
		return nil, err
	}

	// Hämtar filinformation från databasen
	var file model.File
	var createdAt string
	var fileData []byte
	err := r.DB.QueryRow(`
		SELECT id, name, size, content_type, created_at, file_data 
		FROM files WHERE id = ?`, id).Scan(
		&file.ID, &file.Name, &file.Size, &file.ContentType, &createdAt, &fileData)

	if err == sql.ErrNoRows {
		errMsg := fmt.Sprintf("file not found with ID: %s", id)
		log.Printf("Download failed: %s", errMsg)
		return nil, fmt.Errorf(errMsg)
	} else if err != nil {
		log.Printf("Download failed: Error fetching file with ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to fetch file: %v", err)
	}

	file.CreatedAt = createdAt
	if fileData != nil {
		encodedFileData := base64.StdEncoding.EncodeToString(fileData)
		file.FileData = &encodedFileData
	}

	// Hämtar metadata för filen
	metaRows, err := r.DB.Query("SELECT key, value FROM metadata WHERE file_id = ?", file.ID)
	if err != nil {
		log.Printf("Error fetching metadata for file ID %s: %v", file.ID, err)
		return nil, fmt.Errorf("failed to fetch metadata: %v", err)
	}
	defer metaRows.Close()

	var metadata []*model.Metadata
	for metaRows.Next() {
		var meta model.Metadata
		if err := metaRows.Scan(&meta.Key, &meta.Value); err != nil {
			log.Printf("Error scanning metadata row: %v", err)
			return nil, fmt.Errorf("failed to scan metadata row: %v", err)
		}
		metadata = append(metadata, &meta)
	}
	file.Metadata = metadata

	log.Printf("Successfully prepared file %s (ID: %s) for download", file.Name, file.ID)
	return &file, nil
}

// Hello implementerar Query.hello
func (r *queryResolver) Hello(ctx context.Context) (string, error) {
	return "Hello from GraphQL!", nil
}

// User implementerar Todo.user
func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	return &model.User{
		ID:   "1", // Exempel-ID eftersom detta är en dummy-implementation
		Name: "Example User",
	}, nil
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }

// Implementera logAction för att logga händelser
func logAction(action string) {
	log.Printf("[ACTION] %s", action)
}
