-- Remove indexes
DROP INDEX IF EXISTS suppliers_district_id_index;
DROP INDEX IF EXISTS suppliers_city_id_index;

-- Remove columns
ALTER TABLE suppliers 
DROP COLUMN IF EXISTS district_id,
DROP COLUMN IF EXISTS city_id;

