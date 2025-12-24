-- Rollback: Drop the sync function and reset expedition_code_seq to start from 1
-- Note: This may cause duplicate key issues if there are existing records
DROP FUNCTION IF EXISTS sync_expedition_code_sequence();
ALTER SEQUENCE expedition_code_seq RESTART WITH 1;
