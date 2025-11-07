-- Restore telp_number and phone_number columns to expeditions table
ALTER TABLE expeditions ADD COLUMN IF NOT EXISTS telp_number VARCHAR(255);
ALTER TABLE expeditions ADD COLUMN IF NOT EXISTS phone_number VARCHAR(255);

