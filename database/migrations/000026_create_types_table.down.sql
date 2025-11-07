-- Drop indexes
DROP INDEX IF EXISTS types_deleted_at_index;
DROP INDEX IF EXISTS types_updated_at_index;
DROP INDEX IF EXISTS types_created_at_index;
DROP INDEX IF EXISTS types_name_index;
DROP INDEX IF EXISTS types_type_code_index;
DROP INDEX IF EXISTS types_subgroup_id_index;
DROP INDEX IF EXISTS types_id_index;

-- Drop table
DROP TABLE IF EXISTS types;

-- Drop function
DROP FUNCTION IF EXISTS generate_type_code();

-- Drop sequence
DROP SEQUENCE IF EXISTS type_code_seq;

