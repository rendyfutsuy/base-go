DROP INDEX IF EXISTS idx_jwt_tokens_refresh_jti;
DROP INDEX IF EXISTS idx_jwt_tokens_access_jti;

ALTER TABLE jwt_tokens DROP COLUMN IF NOT EXISTS is_used;
ALTER TABLE jwt_tokens DROP COLUMN IF NOT EXISTS refresh_expires_at;
ALTER TABLE jwt_tokens DROP COLUMN IF NOT EXISTS refresh_jti;
ALTER TABLE jwt_tokens DROP COLUMN IF NOT EXISTS access_jti;
