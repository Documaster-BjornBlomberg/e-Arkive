-- database: e-Arkive.db

-- Drop existing tables if they exist
DROP TABLE IF EXISTS metadata;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS nodes;
DROP TABLE IF EXISTS user_settings;
DROP TABLE IF EXISTS users;

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

-- Create table for storing files with BLOB support and node relationship
CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    size INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    created_at TEXT NOT NULL,
    file_data BLOB,
    node_id INTEGER DEFAULT 1,
    FOREIGN KEY (node_id) REFERENCES nodes (id)
);

-- Create index for faster node-based file queries
CREATE INDEX idx_files_node_id ON files(node_id);

-- Create table for storing metadata associated with files
CREATE TABLE IF NOT EXISTS metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (file_id) REFERENCES files (id) ON DELETE CASCADE
);

-- Create users table for authentication
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TEXT NOT NULL
);

-- Create user settings table
CREATE TABLE IF NOT EXISTS user_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create index for faster user settings queries
CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);
-- Create unique index on user_id and key to ensure no duplicate settings
CREATE UNIQUE INDEX idx_user_settings_unique ON user_settings(user_id, key);

-- Insert default admin user with password "admin" (using bcrypt hash)
INSERT INTO users (username, password_hash, created_at)
VALUES ('admin', '$2a$10$G/Yn1SqchSCfNYdN6.LYBemwy8pwMAFwFb30il2wzmkb57wgS2f6q', datetime('now'));