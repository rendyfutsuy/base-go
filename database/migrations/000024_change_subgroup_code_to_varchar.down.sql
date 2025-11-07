-- Rollback: Change subgroup_code back from VARCHAR to BIGINT

-- Step 1: Drop trigger
DROP TRIGGER IF EXISTS trigger_generate_subgroup_code ON sub_groups;

-- Step 2: Drop function
DROP FUNCTION IF EXISTS generate_subgroup_code();

-- Step 3: Convert formatted strings back to integers
-- Remove leading zeros and convert to integer
UPDATE sub_groups 
SET subgroup_code = CAST(CAST(subgroup_code AS BIGINT) AS TEXT)
WHERE subgroup_code IS NOT NULL AND subgroup_code ~ '^[0-9]+$';

-- Step 4: Alter column type back to BIGINT
ALTER TABLE sub_groups 
ALTER COLUMN subgroup_code TYPE BIGINT USING CAST(subgroup_code AS BIGINT);

-- Step 5: Reset DEFAULT to use sequence directly (for rollback)
ALTER TABLE sub_groups 
ALTER COLUMN subgroup_code SET DEFAULT nextval('subgroup_code_seq');

