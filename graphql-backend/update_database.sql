-- database: e-Arkive.db

-- Drop existing tables if they exist
DROP TABLE IF EXISTS metadata;
DROP TABLE IF EXISTS files;

-- Create table for storing files with BLOB support
CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    size INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    created_at TEXT NOT NULL,
    file_data BLOB
);

-- Create table for storing metadata associated with files
CREATE TABLE IF NOT EXISTS metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (file_id) REFERENCES files (id) ON DELETE CASCADE
);