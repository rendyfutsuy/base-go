-- Enable pg_trgm extension for similarity and trigram indexes
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- City and Subdistrict names (search via EXISTS subqueries)
CREATE INDEX IF NOT EXISTS city_name_trgm_idx ON city USING gin (LOWER(REPLACE(name, ' ', '')) gin_trgm_ops);
CREATE INDEX IF NOT EXISTS subdistrict_name_trgm_idx ON subdistrict USING gin (LOWER(REPLACE(name, ' ', '')) gin_trgm_ops);

-- Goods Group (search: name, group_code)
CREATE INDEX IF NOT EXISTS goods_group_name_trgm_idx ON goods_group USING gin (LOWER(REPLACE(name, ' ', '')) gin_trgm_ops);
CREATE INDEX IF NOT EXISTS goods_group_group_code_trgm_idx ON goods_group USING gin (LOWER(REPLACE(group_code, ' ', '')) gin_trgm_ops);

-- Sub Groups (search via goods and other modules by name)
CREATE INDEX IF NOT EXISTS sub_groups_name_trgm_idx ON sub_groups USING gin (LOWER(REPLACE(name, ' ', '')) gin_trgm_ops);

-- Types and Backings (search via goods)
CREATE INDEX IF NOT EXISTS types_name_trgm_idx ON types USING gin (LOWER(REPLACE(name, ' ', '')) gin_trgm_ops);
CREATE INDEX IF NOT EXISTS backings_name_trgm_idx ON backings USING gin (LOWER(REPLACE(name, ' ', '')) gin_trgm_ops);

-- Parameters (search via multiple modules by parameter name)
CREATE INDEX IF NOT EXISTS parameter_name_trgm_idx ON parameter USING gin (LOWER(REPLACE(name, ' ', '')) gin_trgm_ops);

