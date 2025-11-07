-- Rollback: Change identity_type back from UUID to VARCHAR

-- Step 1: Drop foreign key constraint
ALTER TABLE suppliers DROP CONSTRAINT IF EXISTS fk_suppliers_identity_type;

-- Step 2: Drop index
DROP INDEX IF EXISTS suppliers_identity_type_index;

-- Step 3: Add temporary VARCHAR column
ALTER TABLE suppliers ADD COLUMN identity_type_varchar VARCHAR(255);

-- Step 4: Migrate data back: Get parameter name/code
UPDATE suppliers s
SET identity_type_varchar = COALESCE(p.code, p.name)
FROM parameter p
WHERE s.identity_type = p.id
  AND p.deleted_at IS NULL;

-- Step 5: Drop UUID column
ALTER TABLE suppliers DROP COLUMN identity_type;

-- Step 6: Rename VARCHAR column to identity_type
ALTER TABLE suppliers RENAME COLUMN identity_type_varchar TO identity_type;

-- Step 7: Set NOT NULL constraint
ALTER TABLE suppliers ALTER COLUMN identity_type SET NOT NULL;

-- Step 8: Recreate index
CREATE INDEX IF NOT EXISTS suppliers_identity_type_index ON suppliers (identity_type);

