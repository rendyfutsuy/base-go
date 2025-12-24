-- Drop indexes
DROP INDEX IF EXISTS parameter_deleted_at_index;
DROP INDEX IF EXISTS parameter_updated_at_index;
DROP INDEX IF EXISTS parameter_created_at_index;
DROP INDEX IF EXISTS parameter_name_index;
DROP INDEX IF EXISTS parameter_code_index;
DROP INDEX IF EXISTS parameter_id_index;

-- Drop table
DROP TABLE IF EXISTS parameters;

