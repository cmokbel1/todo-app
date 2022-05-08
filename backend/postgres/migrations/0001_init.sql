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
    created_at TIMESTAMPTZ           NOT NULL,
    updated_at TIMESTAMPTZ           NOT NULL
);

CREATE INDEX users_name_idx ON users (name);
CREATE UNIQUE INDEX users_name_key ON users (LOWER(name));
CREATE INDEX users_email_idx ON users (email);
CREATE UNIQUE INDEX users_email_key ON users (LOWER(email));

-- user_credentials represent the list of app user credentials (when auth.source == "app")
CREATE TABLE IF NOT EXISTS user_credentials
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    user_id    BIGINT REFERENCES users (id) ON DELETE CASCADE,
    name       TEXT                  NOT NULL,
    password   TEXT                  NOT NULL,
    -- api_key can be used to auth in lieu of a name+password combination
    api_key    TEXT                  NOT NULL,

    created_at TIMESTAMPTZ           NOT NULL,
    updated_at TIMESTAMPTZ           NOT NULL
);

CREATE INDEX user_app_credentials_name_idx ON user_credentials (name);
CREATE UNIQUE INDEX user_app_credentials_name_key ON user_credentials (LOWER(name));
CREATE UNIQUE INDEX users_api_key_idx ON user_credentials (api_key);

CREATE TABLE IF NOT EXISTS auths
(
    id            BIGSERIAL PRIMARY KEY NOT NULL,
    user_id       BIGINT REFERENCES users (id) ON DELETE CASCADE,
    source        TEXT                  NOT NULL, -- the source of the auth, one of: "app", "github"
    source_id     TEXT                  NOT NULL, -- the id of the user at the source

    access_token  TEXT                  NOT NULL, -- only set when auth created by OAuth
    refresh_token TEXT                  NOT NULL, -- only set when auth created by OAuth

    expiry        TIMESTAMPTZ           NOT NULL, -- the time at which this auth context should be considered invalid
    created_at    TIMESTAMPTZ           NOT NULL,
    updated_at    TIMESTAMPTZ           NOT NULL
);

CREATE UNIQUE INDEX auths_user_id_source_key ON auths (user_id, source); -- one source per user
CREATE UNIQUE INDEX auths_source_source_id_key ON auths (source, source_id); -- one auth per source user

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
DROP TABLE IF EXISTS user_app_credentials;
DROP TABLE IF EXISTS auths;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sessions;