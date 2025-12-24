-- Drop area_code column from city table
ALTER TABLE cities
    DROP COLUMN IF EXISTS area_code;
