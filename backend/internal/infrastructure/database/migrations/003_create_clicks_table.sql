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
    device VARCHAR(50),
    browser VARCHAR(50),
    os VARCHAR(50),
    clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes for clicks table
CREATE INDEX IF NOT EXISTS idx_clicks_short_url_id ON clicks(short_url_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_clicks_short_url_id_clicked_at ON clicks(short_url_id, clicked_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_clicks_ip_address ON clicks(ip_address) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_clicks_country_clicked_at ON clicks(country, clicked_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at ON clicks(clicked_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_clicks_deleted_at ON clicks(deleted_at);

-- Add trigger to update updated_at timestamp
CREATE TRIGGER update_clicks_updated_at 
    BEFORE UPDATE ON clicks 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add trigger to increment click_count on short_urls
CREATE OR REPLACE FUNCTION increment_click_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE short_urls 
    SET click_count = click_count + 1 
    WHERE id = NEW.short_url_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER increment_url_click_count 
    AFTER INSERT ON clicks 
    FOR EACH ROW 
    EXECUTE FUNCTION increment_click_count();