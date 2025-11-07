-- Add back telp_number and phone_number columns to suppliers table
ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS telp_number VARCHAR(255);
ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS phone_number VARCHAR(255);

