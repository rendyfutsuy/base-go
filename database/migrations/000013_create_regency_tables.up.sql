-- Create table province
CREATE TABLE IF NOT EXISTS provinces (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  name VARCHAR(100) NOT NULL UNIQUE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Indexes for province
CREATE INDEX IF NOT EXISTS province_id_index ON provinces (id);
CREATE INDEX IF NOT EXISTS province_name_index ON provinces (name);
CREATE INDEX IF NOT EXISTS province_created_at_index ON provinces (created_at);
CREATE INDEX IF NOT EXISTS province_updated_at_index ON provinces (updated_at);
CREATE INDEX IF NOT EXISTS province_deleted_at_index ON provinces (deleted_at);

-- Create table city
CREATE TABLE IF NOT EXISTS cities (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  province_id UUID NOT NULL REFERENCES provinces(id),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Indexes for city
CREATE INDEX IF NOT EXISTS city_id_index ON cities (id);
CREATE INDEX IF NOT EXISTS city_province_id_index ON cities (province_id);
CREATE INDEX IF NOT EXISTS city_name_index ON cities (name);
CREATE INDEX IF NOT EXISTS city_created_at_index ON cities (created_at);
CREATE INDEX IF NOT EXISTS city_updated_at_index ON cities (updated_at);
CREATE INDEX IF NOT EXISTS city_deleted_at_index ON cities (deleted_at);

-- Unique constraint: city name must be unique within a province
CREATE UNIQUE INDEX IF NOT EXISTS city_in_province ON cities (province_id, name) WHERE deleted_at IS NULL;

-- Create table district
CREATE TABLE IF NOT EXISTS districts (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  city_id UUID NOT NULL REFERENCES cities(id),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Indexes for district
CREATE INDEX IF NOT EXISTS district_id_index ON districts (id);
CREATE INDEX IF NOT EXISTS district_city_id_index ON districts (city_id);
CREATE INDEX IF NOT EXISTS district_name_index ON districts (name);
CREATE INDEX IF NOT EXISTS district_created_at_index ON districts (created_at);
CREATE INDEX IF NOT EXISTS district_updated_at_index ON districts (updated_at);
CREATE INDEX IF NOT EXISTS district_deleted_at_index ON districts (deleted_at);

-- Unique constraint: district name must be unique within a city
CREATE UNIQUE INDEX IF NOT EXISTS district_in_city ON districts (city_id, name) WHERE deleted_at IS NULL;

-- Create table subdistrict
CREATE TABLE IF NOT EXISTS subdistricts (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  district_id UUID NOT NULL REFERENCES districts(id),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Indexes for subdistrict
CREATE INDEX IF NOT EXISTS subdistrict_id_index ON subdistricts (id);
CREATE INDEX IF NOT EXISTS subdistrict_district_id_index ON subdistricts (district_id);
CREATE INDEX IF NOT EXISTS subdistrict_name_index ON subdistricts (name);
CREATE INDEX IF NOT EXISTS subdistrict_created_at_index ON subdistricts (created_at);
CREATE INDEX IF NOT EXISTS subdistrict_updated_at_index ON subdistricts (updated_at);
CREATE INDEX IF NOT EXISTS subdistrict_deleted_at_index ON subdistricts (deleted_at);

-- Unique constraint: subdistrict name must be unique within a district
CREATE UNIQUE INDEX IF NOT EXISTS subdistrict_in_district ON subdistricts (district_id, name) WHERE deleted_at IS NULL;

