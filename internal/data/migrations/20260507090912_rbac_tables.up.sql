-- 1. Roles table with self-referencing parent_id for hierarchy
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    parent_id BIGINT REFERENCES roles(id) ON DELETE SET NULL,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- 2. Groups table
CREATE TABLE groups (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- 3. User-Role Junction (Direct Assignment)
CREATE TABLE user_roles (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    PRIMARY KEY (user_id, role_id)
);

-- 4. User-Group Junction
CREATE TABLE user_groups (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    PRIMARY KEY (user_id, group_id)
);

-- 5. Group-Role Junction (The "G" in GBAC)
CREATE TABLE group_roles (
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    PRIMARY KEY (group_id, role_id)
);

-- 6. Permissions Table
CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE, -- e.g., 'orders:create', 'users:view'
    description TEXT,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- 7. Role-Permission Junction
CREATE TABLE role_permissions (
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    status VARCHAR(1) NOT NULL DEFAULT 'A',
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    PRIMARY KEY (role_id, permission_id)
);

-- 8. Initial Inserts
-- Create Hierarchy: super_admin -> admin -> editor -> viewer
INSERT INTO roles (name, description, parent_id) VALUES 
('viewer', 'Can only read data', NULL),
('editor', 'Can edit data', 1), -- Editor inherits Viewer
('admin', 'Full access to module', 2), -- Admin inherits Editor
('super_admin', 'God mode', 3); -- Super Admin inherits everything

-- Create a Group
INSERT INTO groups (name, description) VALUES 
('Management', 'Users with elevated department privileges');

-- Link Group to a Role
INSERT INTO group_roles (group_id, role_id) VALUES 
(1, 3); -- Management group gets 'admin' role

-- 9. Seed some permissions
INSERT INTO permissions (slug, description) VALUES 
('users:view', 'Can view user profiles'),
('users:edit', 'Can modify user data'),
('orders:view', 'Can view order list'),
('orders:create', 'Can place new orders'),
('orders:delete', 'Can cancel/delete orders');

-- 9. Link Permissions to Roles
-- Let's give 'admin' (ID 3) everything
INSERT INTO role_permissions (role_id, permission_id) 
SELECT 3, id FROM permissions;

-- Let's give 'viewer' (ID 1) only view perms
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions WHERE slug LIKE '%:view';