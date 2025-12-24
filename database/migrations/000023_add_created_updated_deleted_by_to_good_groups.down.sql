-- Rollback: Remove created_by, updated_by, and deleted_by columns from groups table

-- Drop indexes first
DROP INDEX IF EXISTS groups_created_by_index;
DROP INDEX IF EXISTS groups_updated_by_index;
DROP INDEX IF EXISTS groups_deleted_by_index;

-- Drop columns
ALTER TABLE groups
DROP COLUMN IF EXISTS created_by,
DROP COLUMN IF EXISTS updated_by,
DROP COLUMN IF EXISTS deleted_by;
