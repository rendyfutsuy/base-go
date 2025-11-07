-- Change identity_type from VARCHAR to UUID with FK to parameters table

-- Step 1: Drop the existing index on identity_type
DROP INDEX IF EXISTS suppliers_identity_type_index;

-- Step 2: Add a temporary column for UUID
ALTER TABLE suppliers ADD COLUMN identity_type_uuid UUID;

-- Step 3: Migrate existing data: Try to find matching parameter by name
-- This assumes existing identity_type values match parameter.code or parameter.name
-- If no match found, identity_type_uuid will be NULL (you may need to handle this manually)
UPDATE suppliers s
SET identity_type_uuid = p.id
FROM parameter p
WHERE s.identity_type = p.code OR s.identity_type = p.name
  AND p.deleted_at IS NULL;

-- Step 4: Drop the old VARCHAR column
ALTER TABLE suppliers DROP COLUMN identity_type;

-- Step 5: Rename the UUID column to identity_type
ALTER TABLE suppliers RENAME COLUMN identity_type_uuid TO identity_type;

-- Step 6: Set NOT NULL constraint
ALTER TABLE suppliers ALTER COLUMN identity_type SET NOT NULL;

-- Step 7: Add foreign key constraint
ALTER TABLE suppliers 
ADD CONSTRAINT fk_suppliers_identity_type 
FOREIGN KEY (identity_type) REFERENCES parameter(id) ON DELETE RESTRICT;

-- Step 8: Recreate index
CREATE INDEX IF NOT EXISTS suppliers_identity_type_index ON suppliers (identity_type);

