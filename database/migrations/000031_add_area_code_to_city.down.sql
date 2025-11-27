-- Drop area_code column from city table
ALTER TABLE city
    DROP COLUMN IF EXISTS area_code;

