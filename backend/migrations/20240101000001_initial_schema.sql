-- +migrate Up
-- Initial schema for URL Shortener service

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create short_urls table
CREATE TABLE IF NOT EXISTS short_urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    custom_alias BOOLEAN DEFAULT FALSE,
    password_hash VARCHAR(255),
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create clicks table
CREATE TABLE IF NOT EXISTS clicks (
    id SERIAL PRIMARY KEY,
    short_url_id INTEGER NOT NULL REFERENCES short_urls(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT,
    referer TEXT,
    country VARCHAR(2),
    region VARCHAR(100),
    city VARCHAR(100),
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create basic indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_short_urls_short_code ON short_urls(short_code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_short_urls_user_id ON short_urls(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_clicks_short_url_id ON clicks(short_url_id);
CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at ON clicks(clicked_at);

-- +migrate Down
-- Drop tables in reverse order to handle foreign key constraints

DROP INDEX IF EXISTS idx_clicks_clicked_at;
DROP INDEX IF EXISTS idx_clicks_short_url_id;
DROP INDEX IF EXISTS idx_short_urls_user_id;
DROP INDEX IF EXISTS idx_short_urls_short_code;
DROP INDEX IF EXISTS idx_users_email;

DROP TABLE IF EXISTS clicks;
DROP TABLE IF EXISTS short_urls;
DROP TABLE IF EXISTS users;