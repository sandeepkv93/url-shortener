-- Migration: Create webhooks and webhook_deliveries tables
-- Up migration

CREATE TYPE webhook_status AS ENUM ('active', 'inactive', 'failed', 'suspended');
CREATE TYPE webhook_delivery_status AS ENUM ('pending', 'success', 'failed', 'retrying', 'abandoned');

-- Create webhooks table
CREATE TABLE webhooks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(2048) NOT NULL,
    events JSONB NOT NULL DEFAULT '[]'::jsonb,
    secret VARCHAR(255) DEFAULT '',
    status webhook_status NOT NULL DEFAULT 'active',
    
    -- Configuration
    max_retries INTEGER NOT NULL DEFAULT 3 CHECK (max_retries >= 0 AND max_retries <= 10),
    timeout_seconds INTEGER NOT NULL DEFAULT 30 CHECK (timeout_seconds >= 1 AND timeout_seconds <= 300),
    retry_backoff_ms INTEGER NOT NULL DEFAULT 1000 CHECK (retry_backoff_ms >= 100),
    
    -- Statistics
    total_deliveries BIGINT NOT NULL DEFAULT 0,
    success_deliveries BIGINT NOT NULL DEFAULT 0,
    failed_deliveries BIGINT NOT NULL DEFAULT 0,
    last_delivery_at TIMESTAMP WITH TIME ZONE,
    last_success_at TIMESTAMP WITH TIME ZONE,
    last_failure_at TIMESTAMP WITH TIME ZONE,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create webhook_deliveries table
CREATE TABLE webhook_deliveries (
    id BIGSERIAL PRIMARY KEY,
    webhook_id BIGINT NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    event_type VARCHAR(255) NOT NULL,
    status webhook_delivery_status NOT NULL DEFAULT 'pending',
    
    -- Request details
    request_url VARCHAR(2048) NOT NULL,
    request_headers JSONB,
    request_body JSONB,
    
    -- Response details
    response_status INTEGER DEFAULT 0,
    response_headers JSONB,
    response_body TEXT DEFAULT '',
    
    -- Timing and retry logic
    duration BIGINT DEFAULT 0, -- Duration in milliseconds
    attempt_count INTEGER NOT NULL DEFAULT 1,
    next_retry_at TIMESTAMP WITH TIME ZONE,
    
    -- Error handling
    error_message TEXT DEFAULT '',
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for webhooks table
CREATE INDEX idx_webhooks_user_id ON webhooks(user_id);
CREATE INDEX idx_webhooks_status ON webhooks(status);
CREATE INDEX idx_webhooks_events ON webhooks USING GIN(events);
CREATE INDEX idx_webhooks_user_status ON webhooks(user_id, status);
CREATE INDEX idx_webhooks_created_at ON webhooks(created_at);

-- Create indexes for webhook_deliveries table
CREATE INDEX idx_webhook_deliveries_webhook_id ON webhook_deliveries(webhook_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status);
CREATE INDEX idx_webhook_deliveries_event_type ON webhook_deliveries(event_type);
CREATE INDEX idx_webhook_deliveries_created_at ON webhook_deliveries(created_at);
CREATE INDEX idx_webhook_deliveries_next_retry ON webhook_deliveries(next_retry_at) WHERE next_retry_at IS NOT NULL;
CREATE INDEX idx_webhook_deliveries_pending_retries ON webhook_deliveries(status, next_retry_at, attempt_count) 
    WHERE status = 'failed' AND attempt_count < 10;

-- Create composite indexes for common queries
CREATE INDEX idx_webhook_deliveries_webhook_status ON webhook_deliveries(webhook_id, status);
CREATE INDEX idx_webhook_deliveries_webhook_created ON webhook_deliveries(webhook_id, created_at DESC);

-- Create trigger function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at columns
CREATE TRIGGER update_webhooks_updated_at 
    BEFORE UPDATE ON webhooks 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_webhook_deliveries_updated_at 
    BEFORE UPDATE ON webhook_deliveries 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to validate webhook events
CREATE OR REPLACE FUNCTION validate_webhook_events(events JSONB)
RETURNS BOOLEAN AS $$
DECLARE
    valid_events TEXT[] := ARRAY[
        'url.created', 'url.updated', 'url.deleted', 'url.clicked', 'url.expired',
        'analytics.threshold', 'analytics.report',
        'user.registered', 'user.updated',
        'system.error', 'system.alert'
    ];
    event_value TEXT;
BEGIN
    -- Check if events is an array
    IF jsonb_typeof(events) != 'array' THEN
        RETURN FALSE;
    END IF;
    
    -- Check if array is not empty
    IF jsonb_array_length(events) = 0 THEN
        RETURN FALSE;
    END IF;
    
    -- Validate each event
    FOR event_value IN SELECT jsonb_array_elements_text(events)
    LOOP
        IF NOT (event_value = ANY(valid_events)) THEN
            RETURN FALSE;
        END IF;
    END LOOP;
    
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- Add check constraint for webhook events validation
ALTER TABLE webhooks ADD CONSTRAINT valid_webhook_events 
    CHECK (validate_webhook_events(events));

-- Create function to clean up old webhook deliveries
CREATE OR REPLACE FUNCTION cleanup_old_webhook_deliveries(days_to_keep INTEGER DEFAULT 90)
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM webhook_deliveries 
    WHERE created_at < NOW() - INTERVAL '1 day' * days_to_keep;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Create a view for webhook statistics
CREATE OR REPLACE VIEW webhook_stats AS
SELECT 
    w.id,
    w.user_id,
    w.name,
    w.url,
    w.status,
    w.total_deliveries,
    w.success_deliveries,
    w.failed_deliveries,
    CASE 
        WHEN w.total_deliveries > 0 
        THEN ROUND((w.success_deliveries::DECIMAL / w.total_deliveries::DECIMAL) * 100, 2)
        ELSE 0
    END AS success_rate_percent,
    w.last_delivery_at,
    w.last_success_at,
    w.last_failure_at,
    COUNT(wd.id) FILTER (WHERE wd.created_at >= NOW() - INTERVAL '24 hours') AS deliveries_last_24h,
    COUNT(wd.id) FILTER (WHERE wd.created_at >= NOW() - INTERVAL '7 days') AS deliveries_last_7d,
    COUNT(wd.id) FILTER (WHERE wd.created_at >= NOW() - INTERVAL '30 days') AS deliveries_last_30d,
    AVG(wd.duration) FILTER (WHERE wd.status = 'success' AND wd.duration > 0) AS avg_response_time_ms
FROM webhooks w
LEFT JOIN webhook_deliveries wd ON w.id = wd.webhook_id
GROUP BY w.id, w.user_id, w.name, w.url, w.status, w.total_deliveries, 
         w.success_deliveries, w.failed_deliveries, w.last_delivery_at, 
         w.last_success_at, w.last_failure_at;

-- Create sample webhook event types lookup (for reference)
CREATE TABLE webhook_event_types (
    event_type VARCHAR(255) PRIMARY KEY,
    category VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Insert webhook event types
INSERT INTO webhook_event_types (event_type, category, description) VALUES
('url.created', 'url', 'Triggered when a new short URL is created'),
('url.updated', 'url', 'Triggered when a short URL is updated'),
('url.deleted', 'url', 'Triggered when a short URL is deleted'),
('url.clicked', 'url', 'Triggered when a short URL is clicked'),
('url.expired', 'url', 'Triggered when a short URL expires'),
('analytics.threshold', 'analytics', 'Triggered when analytics metrics reach a threshold'),
('analytics.report', 'analytics', 'Triggered when analytics reports are generated'),
('user.registered', 'user', 'Triggered when a new user registers'),
('user.updated', 'user', 'Triggered when user profile is updated'),
('system.error', 'system', 'Triggered when system errors occur'),
('system.alert', 'system', 'Triggered when system alerts are raised');

-- Add comments for documentation
COMMENT ON TABLE webhooks IS 'User-configured webhooks for receiving event notifications';
COMMENT ON TABLE webhook_deliveries IS 'Webhook delivery attempts and their results';
COMMENT ON TABLE webhook_event_types IS 'Reference table for available webhook event types';

COMMENT ON COLUMN webhooks.events IS 'JSON array of webhook events this webhook subscribes to';
COMMENT ON COLUMN webhooks.secret IS 'Secret key used for webhook signature verification';
COMMENT ON COLUMN webhooks.max_retries IS 'Maximum number of retry attempts for failed deliveries';
COMMENT ON COLUMN webhooks.timeout_seconds IS 'HTTP timeout for webhook requests in seconds';
COMMENT ON COLUMN webhook_deliveries.duration IS 'Request duration in milliseconds';
COMMENT ON COLUMN webhook_deliveries.attempt_count IS 'Number of delivery attempts made';
COMMENT ON COLUMN webhook_deliveries.next_retry_at IS 'Scheduled time for next retry attempt';

-- Down migration (commented out for safety)
/*
DROP VIEW IF EXISTS webhook_stats;
DROP FUNCTION IF EXISTS cleanup_old_webhook_deliveries(INTEGER);
DROP FUNCTION IF EXISTS validate_webhook_events(JSONB);
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
DROP TABLE IF EXISTS webhook_event_types;
DROP TABLE IF EXISTS webhook_deliveries;
DROP TABLE IF EXISTS webhooks;
DROP TYPE IF EXISTS webhook_delivery_status;
DROP TYPE IF EXISTS webhook_status;
*/