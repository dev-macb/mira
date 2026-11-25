-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'admin' CHECK (role = 'admin'),
    active INTEGER NOT NULL DEFAULT 1 CHECK (active IN (0, 1)),
    last_sign_in_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id),
    token_hash TEXT NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    revoked INTEGER NOT NULL DEFAULT 0 CHECK (revoked IN (0, 1))
);
CREATE INDEX refresh_tokens_user ON refresh_tokens (user_id);

CREATE TABLE projects (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    client TEXT,
    domain TEXT,
    description TEXT,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE deployments (
    id INTEGER PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects (id),
    revision_sha TEXT NOT NULL,
    services TEXT NOT NULL,
    workflow_run_id INTEGER,
    workflow_run_url TEXT,
    status TEXT NOT NULL,
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME,
    log TEXT NOT NULL DEFAULT ''
);
CREATE INDEX deployments_project_started ON deployments (project_id, started_at DESC);
CREATE UNIQUE INDEX deployments_active_project ON deployments (project_id) WHERE finished_at IS NULL;

CREATE TABLE audit_log (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id),
    action TEXT NOT NULL,
    target TEXT NOT NULL,
    result TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
