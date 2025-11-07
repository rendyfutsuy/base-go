-- Create table province
CREATE TABLE IF NOT EXISTS province (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  name VARCHAR(100) NOT NULL UNIQUE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Indexes for province
CREATE INDEX IF NOT EXISTS province_id_index ON province (id);
CREATE INDEX IF NOT EXISTS province_name_index ON province (name);
CREATE INDEX IF NOT EXISTS province_created_at_index ON province (created_at);
CREATE INDEX IF NOT EXISTS province_updated_at_index ON province (updated_at);
CREATE INDEX IF NOT EXISTS province_deleted_at_index ON province (deleted_at);

-- Create table city
CREATE TABLE IF NOT EXISTS city (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  province_id UUID NOT NULL REFERENCES province(id),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Indexes for city
CREATE INDEX IF NOT EXISTS city_id_index ON city (id);
CREATE INDEX IF NOT EXISTS city_province_id_index ON city (province_id);
CREATE INDEX IF NOT EXISTS city_name_index ON city (name);
CREATE INDEX IF NOT EXISTS city_created_at_index ON city (created_at);
CREATE INDEX IF NOT EXISTS city_updated_at_index ON city (updated_at);
CREATE INDEX IF NOT EXISTS city_deleted_at_index ON city (deleted_at);

-- Unique constraint: city name must be unique within a province
CREATE UNIQUE INDEX IF NOT EXISTS city_in_province ON city (province_id, name) WHERE deleted_at IS NULL;

-- Create table district
CREATE TABLE IF NOT EXISTS district (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  city_id UUID NOT NULL REFERENCES city(id),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Indexes for district
CREATE INDEX IF NOT EXISTS district_id_index ON district (id);
CREATE INDEX IF NOT EXISTS district_city_id_index ON district (city_id);
CREATE INDEX IF NOT EXISTS district_name_index ON district (name);
CREATE INDEX IF NOT EXISTS district_created_at_index ON district (created_at);
CREATE INDEX IF NOT EXISTS district_updated_at_index ON district (updated_at);
CREATE INDEX IF NOT EXISTS district_deleted_at_index ON district (deleted_at);

-- Unique constraint: district name must be unique within a city
CREATE UNIQUE INDEX IF NOT EXISTS district_in_city ON district (city_id, name) WHERE deleted_at IS NULL;

-- Create table subdistrict
CREATE TABLE IF NOT EXISTS subdistrict (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  district_id UUID NOT NULL REFERENCES district(id),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Indexes for subdistrict
CREATE INDEX IF NOT EXISTS subdistrict_id_index ON subdistrict (id);
CREATE INDEX IF NOT EXISTS subdistrict_district_id_index ON subdistrict (district_id);
CREATE INDEX IF NOT EXISTS subdistrict_name_index ON subdistrict (name);
CREATE INDEX IF NOT EXISTS subdistrict_created_at_index ON subdistrict (created_at);
CREATE INDEX IF NOT EXISTS subdistrict_updated_at_index ON subdistrict (updated_at);
CREATE INDEX IF NOT EXISTS subdistrict_deleted_at_index ON subdistrict (deleted_at);

-- Unique constraint: subdistrict name must be unique within a district
CREATE UNIQUE INDEX IF NOT EXISTS subdistrict_in_district ON subdistrict (district_id, name) WHERE deleted_at IS NULL;

