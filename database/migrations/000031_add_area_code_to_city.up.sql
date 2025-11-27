-- Add area_code column to city table
ALTER TABLE city
    ADD COLUMN IF NOT EXISTS area_code VARCHAR(50);

COMMENT ON COLUMN city.area_code IS 'City dialing area code';

