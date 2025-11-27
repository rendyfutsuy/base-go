-- Rollback trigram search indexes and extension

DROP INDEX IF EXISTS parameter_name_trgm_idx;

DROP INDEX IF EXISTS goods_goods_name_trgm_idx;
DROP INDEX IF EXISTS goods_full_code_trgm_idx;

DROP INDEX IF EXISTS backings_name_trgm_idx;
DROP INDEX IF EXISTS types_name_trgm_idx;

DROP INDEX IF EXISTS sub_groups_name_trgm_idx;

DROP INDEX IF EXISTS goods_group_group_code_trgm_idx;
DROP INDEX IF EXISTS goods_group_name_trgm_idx;

DROP INDEX IF EXISTS subdistrict_name_trgm_idx;
DROP INDEX IF EXISTS city_name_trgm_idx;

-- Finally, drop extension (optional)
DROP EXTENSION IF EXISTS pg_trgm;

