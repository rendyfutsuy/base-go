-- Add district_id and city_id columns to suppliers table
ALTER TABLE suppliers 
ADD COLUMN IF NOT EXISTS district_id UUID REFERENCES district(id),
ADD COLUMN IF NOT EXISTS city_id UUID REFERENCES city(id);

-- Create indexes for new columns
CREATE INDEX IF NOT EXISTS suppliers_district_id_index ON suppliers (district_id);
CREATE INDEX IF NOT EXISTS suppliers_city_id_index ON suppliers (city_id);

