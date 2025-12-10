CREATE TABLE users (
    id UUID PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE licenses (
    id UUID PRIMARY KEY,
    owner TEXT NOT NULL,
    max_users INT,
    max_messages INT,
    issued_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    signature TEXT
);