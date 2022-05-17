-- +goose Up
CREATE TABLE sessions
(
    token  TEXT PRIMARY KEY,
    data   BYTEA       NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE TABLE IF NOT EXISTS users
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    name       TEXT                  NOT NULL,
    email      TEXT,
    -- api_key can be used to auth in lieu of an identity provider
    api_key    TEXT                  NOT NULL,
    created_at TIMESTAMPTZ           NOT NULL,
    updated_at TIMESTAMPTZ           NOT NULL
);

CREATE INDEX users_name_idx ON users (name);
CREATE UNIQUE INDEX users_name_key ON users (LOWER(name));
CREATE INDEX users_email_idx ON users (email);
CREATE UNIQUE INDEX users_email_key ON users (LOWER(email));
CREATE UNIQUE INDEX users_apikey_key ON users (api_key);

CREATE TABLE IF NOT EXISTS lists
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    -- user_id is the user who created the list
    user_id    BIGINT REFERENCES users (id) ON DELETE CASCADE,
    name       TEXT                  NOT NULL,
    completed  BOOLEAN               NOT NULL,
    created_at TIMESTAMPTZ           NOT NULL,
    updated_at TIMESTAMPTZ           NOT NULL
);

CREATE TABLE IF NOT EXISTS items
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    user_id    BIGINT REFERENCES users (id) ON DELETE CASCADE,
    list_id    BIGINT REFERENCES lists (id) ON DELETE CASCADE,
    name       TEXT                  NOT NULL,
    completed  BOOLEAN               NOT NULL,
    created_at TIMESTAMPTZ           NOT NULL,
    updated_at TIMESTAMPTZ           NOT NULL
);

-- +goose Down

DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS lists;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sessions;