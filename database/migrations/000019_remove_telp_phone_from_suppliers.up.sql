-- Remove telp_number and phone_number columns from suppliers table
ALTER TABLE suppliers DROP COLUMN IF EXISTS telp_number;
ALTER TABLE suppliers DROP COLUMN IF EXISTS phone_number;

