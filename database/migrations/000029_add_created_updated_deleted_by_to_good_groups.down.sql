-- Rollback: Remove created_by, updated_by, and deleted_by columns from goods_group table

-- Drop indexes first
DROP INDEX IF EXISTS goods_group_created_by_index;
DROP INDEX IF EXISTS goods_group_updated_by_index;
DROP INDEX IF EXISTS goods_group_deleted_by_index;

-- Drop columns
ALTER TABLE goods_group
DROP COLUMN IF EXISTS created_by,
DROP COLUMN IF EXISTS updated_by,
DROP COLUMN IF EXISTS deleted_by;

