-- Create short_urls table
CREATE TABLE IF NOT EXISTS short_urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    custom_alias BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    click_count BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes for short_urls table
CREATE UNIQUE INDEX IF NOT EXISTS idx_short_urls_short_code ON short_urls(short_code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_short_urls_user_id ON short_urls(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_short_urls_user_id_created_at ON short_urls(user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_short_urls_short_code_active ON short_urls(short_code, is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_short_urls_expires_at ON short_urls(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_short_urls_deleted_at ON short_urls(deleted_at);

-- Add trigger to update updated_at timestamp
CREATE TRIGGER update_short_urls_updated_at 
    BEFORE UPDATE ON short_urls 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();