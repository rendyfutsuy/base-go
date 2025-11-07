-- Change subgroup_code from BIGINT to VARCHAR(50) with formatted 2-digit code

-- Step 1: Convert existing numeric codes to formatted strings (while still BIGINT column)
-- Update existing data to match the new format by converting to text first
-- We'll need to alter column first to allow text operations, but we'll do conversion via a temporary approach
-- Actually, we need to alter first to VARCHAR, then format existing data

-- Step 1: Alter column type from BIGINT to VARCHAR(50) (temporarily, will be formatted later)
ALTER TABLE sub_groups 
ALTER COLUMN subgroup_code TYPE VARCHAR(50) USING subgroup_code::TEXT;

-- Step 2: Convert existing codes to formatted 2-digit strings
-- Format all existing codes to 2-digit format (e.g., "1" -> "01", "5" -> "05", "10" stays "10")
UPDATE sub_groups 
SET subgroup_code = LPAD(subgroup_code, 2, '0')
WHERE subgroup_code IS NOT NULL
  AND subgroup_code ~ '^[0-9]+$';

-- Step 3: Create a function to generate formatted subgroup_code for future inserts
CREATE OR REPLACE FUNCTION generate_subgroup_code()
RETURNS TRIGGER AS $$
DECLARE
  next_val BIGINT;
  formatted_code VARCHAR(50);
BEGIN
  -- Only generate if subgroup_code is NULL, empty string, or contains only whitespace
  -- This handles cases where GORM might send empty string instead of NULL
  IF NEW.subgroup_code IS NULL OR TRIM(NEW.subgroup_code) = '' THEN
    -- Get next value from sequence
    next_val := nextval('subgroup_code_seq');
    
    -- Format as 2-digit string with leading zero (e.g., 1 -> "01", 5 -> "05", 10 -> "10")
    formatted_code := LPAD(next_val::TEXT, 2, '0');
    
    -- Set the formatted code
    NEW.subgroup_code := formatted_code;
  END IF;
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Step 4: Remove DEFAULT value from column (we'll use trigger instead)
ALTER TABLE sub_groups 
ALTER COLUMN subgroup_code DROP DEFAULT;

-- Step 5: Create trigger to auto-generate formatted code on INSERT
-- Trigger will always run, function will check if code needs to be generated
DROP TRIGGER IF EXISTS trigger_generate_subgroup_code ON sub_groups;
CREATE TRIGGER trigger_generate_subgroup_code
  BEFORE INSERT ON sub_groups
  FOR EACH ROW
  EXECUTE FUNCTION generate_subgroup_code();

-- Step 6: Reset sequence to match current max value (if needed)
-- This ensures the sequence continues correctly after conversion
-- Note: After converting to VARCHAR, we need to parse the numeric value from string
DO $$
DECLARE
  max_code_val BIGINT;
BEGIN
  -- Get the maximum numeric value from existing codes (parse from formatted string)
  -- Since codes are now formatted as "01", "02", etc., we parse them
  SELECT COALESCE(MAX(CAST(subgroup_code AS BIGINT)), 0)
  INTO max_code_val
  FROM sub_groups
  WHERE deleted_at IS NULL
    AND subgroup_code IS NOT NULL
    AND subgroup_code ~ '^[0-9]+$'; -- Only numeric codes in string format

  -- Set the sequence to start from max_code_val + 1
  PERFORM setval('subgroup_code_seq', max_code_val + 1, false);
END $$;

