CREATE TABLE IF NOT EXISTS "order" {
    id UUID PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
}