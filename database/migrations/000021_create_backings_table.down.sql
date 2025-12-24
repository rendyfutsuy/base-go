-- Drop indexes
DROP INDEX IF EXISTS backings_deleted_at_index;
DROP INDEX IF EXISTS backings_updated_at_index;
DROP INDEX IF EXISTS backings_created_at_index;
DROP INDEX IF EXISTS backings_name_index;
DROP INDEX IF EXISTS backings_backing_code_index;
DROP INDEX IF EXISTS backings_type_id_index;
DROP INDEX IF EXISTS backings_id_index;
DROP INDEX IF EXISTS backing_in_type;

-- Drop table
DROP TABLE IF EXISTS backings;

-- Drop function
DROP FUNCTION IF EXISTS generate_backing_code();

-- Drop sequence
DROP SEQUENCE IF EXISTS backing_code_seq;
