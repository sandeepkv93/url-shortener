-- +migrate Up
-- Add analytics improvements to track more detailed click information

-- Add additional fields to clicks table for better analytics
ALTER TABLE clicks ADD COLUMN IF NOT EXISTS device_type VARCHAR(50);
ALTER TABLE clicks ADD COLUMN IF NOT EXISTS browser VARCHAR(100);
ALTER TABLE clicks ADD COLUMN IF NOT EXISTS os VARCHAR(100);
ALTER TABLE clicks ADD COLUMN IF NOT EXISTS is_mobile BOOLEAN DEFAULT FALSE;
ALTER TABLE clicks ADD COLUMN IF NOT EXISTS is_bot BOOLEAN DEFAULT FALSE;
ALTER TABLE clicks ADD COLUMN IF NOT EXISTS utm_source VARCHAR(255);
ALTER TABLE clicks ADD COLUMN IF NOT EXISTS utm_medium VARCHAR(255);
ALTER TABLE clicks ADD COLUMN IF NOT EXISTS utm_campaign VARCHAR(255);

-- Add performance indexes for analytics queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clicks_country_clicked_at ON clicks(country, clicked_at DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clicks_device_type ON clicks(device_type);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clicks_is_mobile ON clicks(is_mobile);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_clicks_utm_source ON clicks(utm_source);

-- Add stats tracking to short_urls table
ALTER TABLE short_urls ADD COLUMN IF NOT EXISTS total_clicks BIGINT DEFAULT 0;
ALTER TABLE short_urls ADD COLUMN IF NOT EXISTS unique_clicks BIGINT DEFAULT 0;
ALTER TABLE short_urls ADD COLUMN IF NOT EXISTS last_clicked_at TIMESTAMP WITH TIME ZONE;

-- Create function to update click statistics
CREATE OR REPLACE FUNCTION update_url_stats()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE short_urls 
    SET 
        total_clicks = total_clicks + 1,
        last_clicked_at = NEW.clicked_at
    WHERE id = NEW.short_url_id;
    
    -- Update unique clicks (simplified - in production would be more sophisticated)
    UPDATE short_urls 
    SET unique_clicks = (
        SELECT COUNT(DISTINCT ip_address) 
        FROM clicks 
        WHERE short_url_id = NEW.short_url_id
    )
    WHERE id = NEW.short_url_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for automatic stats updates
CREATE TRIGGER trigger_update_url_stats
    AFTER INSERT ON clicks
    FOR EACH ROW
    EXECUTE FUNCTION update_url_stats();

-- +migrate Down
-- Remove analytics improvements

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_url_stats ON clicks;
DROP FUNCTION IF EXISTS update_url_stats();

-- Remove stats columns from short_urls
ALTER TABLE short_urls DROP COLUMN IF EXISTS last_clicked_at;
ALTER TABLE short_urls DROP COLUMN IF EXISTS unique_clicks;
ALTER TABLE short_urls DROP COLUMN IF EXISTS total_clicks;

-- Drop analytics indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_clicks_utm_source;
DROP INDEX CONCURRENTLY IF EXISTS idx_clicks_is_mobile;
DROP INDEX CONCURRENTLY IF EXISTS idx_clicks_device_type;
DROP INDEX CONCURRENTLY IF EXISTS idx_clicks_country_clicked_at;

-- Remove analytics columns from clicks table
ALTER TABLE clicks DROP COLUMN IF EXISTS utm_campaign;
ALTER TABLE clicks DROP COLUMN IF EXISTS utm_medium;
ALTER TABLE clicks DROP COLUMN IF EXISTS utm_source;
ALTER TABLE clicks DROP COLUMN IF EXISTS is_bot;
ALTER TABLE clicks DROP COLUMN IF EXISTS is_mobile;
ALTER TABLE clicks DROP COLUMN IF EXISTS os;
ALTER TABLE clicks DROP COLUMN IF EXISTS browser;
ALTER TABLE clicks DROP COLUMN IF EXISTS device_type;