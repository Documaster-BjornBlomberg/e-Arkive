-- Update database schema to support user permissions system

-- Create groups table if it doesn't exist
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create group_members junction table for many-to-many relationship between users and groups
CREATE TABLE IF NOT EXISTS group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    is_admin BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (group_id, user_id)
);

-- Add permission fields to nodes table if they don't exist
ALTER TABLE nodes ADD COLUMN IF NOT EXISTS owner_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE nodes ADD COLUMN IF NOT EXISTS owner_group_id INTEGER REFERENCES groups(id) ON DELETE SET NULL;
ALTER TABLE nodes ADD COLUMN IF NOT EXISTS permissions INTEGER NOT NULL DEFAULT 7; -- Default: Read + Modify + Delete

-- Create an index for faster permission lookups
CREATE INDEX IF NOT EXISTS idx_nodes_permissions ON nodes(owner_user_id, owner_group_id, permissions);
CREATE INDEX IF NOT EXISTS idx_nodes_owner_user ON nodes (owner_user_id);
CREATE INDEX IF NOT EXISTS idx_nodes_owner_group ON nodes (owner_group_id);

-- Insert the default admin group
INSERT OR IGNORE INTO groups (id, name) VALUES (1, 'Administrators');

-- Add the first user (ID 1) to the admin group if they exist
INSERT OR IGNORE INTO group_members (user_id, group_id)
SELECT 1, 1 
WHERE EXISTS (SELECT 1 FROM users WHERE id = 1);

-- Set default owner and permissions for existing nodes
UPDATE nodes 
SET owner_user_id = (SELECT id FROM users ORDER BY id LIMIT 1), 
    permissions = 7 -- Grant all permissions (read + write + delete)
WHERE owner_user_id IS NULL;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON group_members (user_id);
CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members (group_id);