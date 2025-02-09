CREATE DATABASE IF NOT EXISTS mini_url;

CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    clicks INT DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_slug ON urls(slug);
