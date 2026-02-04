ALTER TABLE jwt_tokens ADD COLUMN IF NOT EXISTS access_jti TEXT;
ALTER TABLE jwt_tokens ADD COLUMN IF NOT EXISTS refresh_jti TEXT;
ALTER TABLE jwt_tokens ADD COLUMN IF NOT EXISTS refresh_expires_at TIMESTAMP;
ALTER TABLE jwt_tokens ADD COLUMN IF NOT EXISTS is_used BOOLEAN DEFAULT FALSE;

-- Optional: Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_jwt_tokens_access_jti ON jwt_tokens(access_jti);
CREATE INDEX IF NOT EXISTS idx_jwt_tokens_refresh_jti ON jwt_tokens(refresh_jti);
