-- database: e-Arkive.db

-- Drop existing tables if they exist
DROP TABLE IF EXISTS metadata;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS nodes;

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

-- Create table for hierarchical nodes structure
CREATE TABLE IF NOT EXISTS nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    parent_id INTEGER,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES nodes (id) ON DELETE RESTRICT
);

-- Index på parent_id för snabbare sökning
CREATE INDEX idx_nodes_parent_id ON nodes(parent_id);

-- Insert Root node as a starting point
INSERT INTO nodes (name, parent_id, created_at, updated_at)
VALUES ('Root', NULL, datetime('now'), datetime('now'));