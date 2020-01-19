-- Modules
CREATE TABLE IF NOT EXISTS modules (
    slug VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL
);

-- Users, organizations and roles
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    password VARCHAR(128) NOT NULL,
    email VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(32),
    lastname VARCHAR(32),
    role VARCHAR(32) NOT NULL,
    validated BOOLEAN DEFAULT FALSE,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY,
    organization_id UUID REFERENCES organizations(id),
    name VARCHAR(32) NOT NULL
);

CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY,
    user_id UUID,
    organization_id UUID,
    role_id UUID,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS permissions (
    module_slug VARCHAR(32),
    role_id  UUID,
    c BOOLEAN DEFAULT FALSE,
    r BOOLEAN DEFAULT FALSE,
    u BOOLEAN DEFAULT FALSE,
    d BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (module_slug, role_id),
    FOREIGN KEY (module_slug) REFERENCES modules(slug),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);