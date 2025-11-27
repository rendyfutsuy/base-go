-- Add created_by, updated_by, and deleted_by columns to goods_group table if they don't exist
-- This migration ensures these columns exist even if they were missing from the initial table creation

ALTER TABLE goods_group 
ADD COLUMN IF NOT EXISTS created_by VARCHAR(255),
ADD COLUMN IF NOT EXISTS updated_by VARCHAR(255),
ADD COLUMN IF NOT EXISTS deleted_by VARCHAR(255);

-- Create indexes for the new columns (if they don't exist)
CREATE INDEX IF NOT EXISTS goods_group_created_by_index ON goods_group (created_by);
CREATE INDEX IF NOT EXISTS goods_group_updated_by_index ON goods_group (updated_by);
CREATE INDEX IF NOT EXISTS goods_group_deleted_by_index ON goods_group (deleted_by);

