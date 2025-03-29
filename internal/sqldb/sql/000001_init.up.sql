CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    roles TEXT[] NOT NULL, 
    password_hash TEXT NOT NULL, 
    enabled BOOLEAN NOT NULL, 
    date_created TIMESTAMP NOT NULL, 
    date_updated TIMESTAMP NOT NULL
);