-- database: e-Arkive.db
-- Kombinerad version av update_database.sql och update_permissions.sql

-- Drop existing tables if they exist
DROP TABLE IF EXISTS metadata;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS nodes;
DROP TABLE IF EXISTS user_settings;
DROP TABLE IF EXISTS users;

-- Create users table for authentication
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    name TEXT,
    password_hash TEXT NOT NULL,
    created_at TEXT NOT NULL
);

-- Create table for user groups
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL
);

-- Create table for group members (relation between users and groups)
CREATE TABLE IF NOT EXISTS group_members (
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);

-- Create index for faster group membership lookup
CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON group_members(user_id);
CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);
-- Create unique constraint to prevent duplicate memberships
CREATE UNIQUE INDEX IF NOT EXISTS idx_group_members_unique ON group_members(user_id, group_id);

-- Create table for hierarchical nodes structure with permission support
CREATE TABLE IF NOT EXISTS nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    parent_id INTEGER,
    owner_user_id INTEGER,
    owner_group_id INTEGER,
    permissions INTEGER DEFAULT 63, -- Default: 111111 binary (all permissions)
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES nodes (id) ON DELETE RESTRICT,
    FOREIGN KEY (owner_user_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (owner_group_id) REFERENCES groups (id) ON DELETE SET NULL
);

-- Index på parent_id för snabbare sökning
CREATE INDEX IF NOT EXISTS idx_nodes_parent_id ON nodes(parent_id);
-- Index för snabbare sökning på ägare
CREATE INDEX IF NOT EXISTS idx_nodes_owner_user_id ON nodes(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_nodes_owner_group_id ON nodes(owner_group_id);

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
CREATE INDEX IF NOT EXISTS idx_files_node_id ON files(node_id);

-- Create table for storing metadata associated with files
CREATE TABLE IF NOT EXISTS metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (file_id) REFERENCES files (id) ON DELETE CASCADE
);

-- Create index for faster user settings queries
CREATE INDEX IF NOT EXISTS idx_user_settings_user_id ON user_settings(user_id);
-- Create unique index on user_id and key to ensure no duplicate settings
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_settings_unique ON user_settings(user_id, key);

-- Create Administrators group if not exists
INSERT OR IGNORE INTO groups (id, name, created_at) 
VALUES (1, 'Administrators', datetime('now'));

-- Insert default admin user with password "admin" (using bcrypt hash)
INSERT OR IGNORE INTO users (id, username, password_hash, created_at) 
VALUES (1, 'admin', '$2a$10$G/Yn1SqchSCfNYdN6.LYBemwy8pwMAFwFb30il2wzmkb57wgS2f6q', datetime('now'));

-- Add admin user to Administrators group
INSERT OR IGNORE INTO group_members (user_id, group_id, created_at)
VALUES (1, 1, datetime('now'));

-- Insert Root node as a starting point (admin group owns root node with full permissions)
INSERT OR IGNORE INTO nodes (id, name, parent_id, owner_user_id, owner_group_id, permissions, created_at, updated_at)
VALUES (1, 'Root', NULL, 1, 1, 63, datetime('now'), datetime('now'));

-- Set ownership of existing nodes to admin and admin group
UPDATE nodes SET owner_user_id = 1, owner_group_id = 1, permissions = 63 WHERE owner_user_id IS NULL;

-- Update existing users table to add name column if it doesn't exist
PRAGMA foreign_keys=off;

-- Check if name column exists, if not add it
BEGIN TRANSACTION;

-- Create a temporary table with the new schema
CREATE TABLE IF NOT EXISTS users_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    name TEXT,
    password_hash TEXT NOT NULL,
    created_at TEXT NOT NULL
);

-- Copy data from the old table to the new one
INSERT INTO users_new (id, username, name, password_hash, created_at)
SELECT id, username, username, password_hash, created_at FROM users;

-- Drop the old table
DROP TABLE users;

-- Rename the new table to the old table name
ALTER TABLE users_new RENAME TO users;

COMMIT;
PRAGMA foreign_keys=on;
