-- Add created_by, updated_by, and deleted_by columns to groups table if they don't exist
-- This migration ensures these columns exist even if they were missing from the initial table creation

ALTER TABLE groups 
ADD COLUMN IF NOT EXISTS created_by VARCHAR(255),
ADD COLUMN IF NOT EXISTS updated_by VARCHAR(255),
ADD COLUMN IF NOT EXISTS deleted_by VARCHAR(255);

-- Create indexes for the new columns (if they don't exist)
CREATE INDEX IF NOT EXISTS groups_created_by_index ON groups (created_by);
CREATE INDEX IF NOT EXISTS groups_updated_by_index ON groups (updated_by);
CREATE INDEX IF NOT EXISTS groups_deleted_by_index ON groups (deleted_by);

