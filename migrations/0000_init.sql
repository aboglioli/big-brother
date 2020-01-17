CREATE DATABASE test;
GRANT ALL PRIVILEGES ON DATABASE test TO admin;

CREATE DATABASE users_and_organizations;
GRANT ALL PRIVILEGES ON DATABASE users_and_organizations TO admin;

\c users_and_organizations
-- Modules
CREATE TABLE IF NOT EXISTS modules (
    id UUID PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    slug VARCHAR(32) NOT NULL
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
    user_id UUID,
    organization_id UUID,
    role_id UUID,
    PRIMARY KEY (user_id, organization_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS permissions (
    module_id UUID,
    role_id  UUID,
    permission SMALLINT NOT NULL,
    PRIMARY KEY (module_id, role_id),
    FOREIGN KEY (module_id) REFERENCES modules(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);