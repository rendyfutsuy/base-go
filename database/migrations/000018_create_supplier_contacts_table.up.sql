-- Create supplier_contacts table
CREATE TABLE IF NOT EXISTS supplier_contacts (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
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
CREATE INDEX IF NOT EXISTS supplier_contacts_id_index ON supplier_contacts (id);
CREATE INDEX IF NOT EXISTS supplier_contacts_supplier_id_index ON supplier_contacts (supplier_id);
CREATE INDEX IF NOT EXISTS supplier_contacts_phone_number_index ON supplier_contacts (phone_number);
CREATE INDEX IF NOT EXISTS supplier_contacts_is_primary_index ON supplier_contacts (is_primary);
CREATE INDEX IF NOT EXISTS supplier_contacts_deleted_at_index ON supplier_contacts (deleted_at);

-- Constraint: Only one primary contact per supplier (excluding soft-deleted records)
-- Note: This uses a partial unique index to allow multiple false values but only one true value per supplier
CREATE UNIQUE INDEX IF NOT EXISTS supplier_contacts_supplier_primary_unique 
ON supplier_contacts (supplier_id) 
WHERE is_primary = true AND deleted_at IS NULL;

