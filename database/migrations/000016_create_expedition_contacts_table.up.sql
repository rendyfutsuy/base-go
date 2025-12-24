-- Create expedition_contacts table
CREATE TABLE IF NOT EXISTS expedition_contacts (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  expedition_id UUID NOT NULL REFERENCES expeditions(id) ON DELETE CASCADE,
  phone_type VARCHAR(50) NOT NULL,
  phone_number VARCHAR(50) NOT NULL UNIQUE,
  is_primary BOOLEAN DEFAULT false,
  created_at TIMESTAMP,
  created_by VARCHAR(255),
  updated_at TIMESTAMP,
  updated_by VARCHAR(255),
  deleted_at TIMESTAMP,
  deleted_by VARCHAR(255)
);

-- Indexes
CREATE INDEX IF NOT EXISTS expedition_contacts_id_index ON expedition_contacts (id);
CREATE INDEX IF NOT EXISTS expedition_contacts_expedition_id_index ON expedition_contacts (expedition_id);
CREATE INDEX IF NOT EXISTS expedition_contacts_phone_number_index ON expedition_contacts (phone_number);
CREATE INDEX IF NOT EXISTS expedition_contacts_is_primary_index ON expedition_contacts (is_primary);
CREATE INDEX IF NOT EXISTS expedition_contacts_deleted_at_index ON expedition_contacts (deleted_at);

-- Constraint: Only one primary contact per expedition (excluding soft-deleted records)
-- Note: This uses a partial unique index to allow multiple false values but only one true value per expedition
CREATE UNIQUE INDEX IF NOT EXISTS expedition_contacts_expedition_primary_unique 
ON expedition_contacts (expedition_id) 
WHERE is_primary = true AND deleted_at IS NULL;
