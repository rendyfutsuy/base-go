-- Sync expedition_code_seq sequence with the maximum expedition_code value in expeditions table
-- This fixes the issue where sequence gets out of sync after rollbacks or failed inserts

-- Create a function to sync the sequence (can be called manually if needed)
CREATE OR REPLACE FUNCTION sync_expedition_code_sequence()
RETURNS void AS $$
DECLARE
  max_code_val BIGINT;
BEGIN
  -- Get the maximum numeric value of expedition_code from expeditions table
  -- Convert VARCHAR expedition_code to BIGINT, handling leading zeros
  -- Filter out NULL or empty strings, and only consider non-deleted records
  SELECT COALESCE(MAX(CAST(expedition_code AS BIGINT)), 0)
  INTO max_code_val
  FROM expeditions
  WHERE deleted_at IS NULL 
    AND expedition_code IS NOT NULL 
    AND expedition_code != ''
    AND expedition_code ~ '^[0-9]+$'; -- Only numeric codes

  -- Set the sequence to start from max_code_val + 1
  -- This ensures the next generated code will be unique
  -- Using false means the next nextval() will return max_code_val + 1
  PERFORM setval('expedition_code_seq', max_code_val + 1, false);
END;
$$ LANGUAGE plpgsql;

-- Execute the sync function immediately
SELECT sync_expedition_code_sequence();

