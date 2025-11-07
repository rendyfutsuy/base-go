-- Remove telp_number and phone_number columns from expeditions table
ALTER TABLE expeditions DROP COLUMN IF EXISTS telp_number;
ALTER TABLE expeditions DROP COLUMN IF EXISTS phone_number;

