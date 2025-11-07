-- Create sequence for expedition_code starting from 1
CREATE SEQUENCE IF NOT EXISTS expedition_code_seq START WITH 1 INCREMENT BY 1;

-- Function to generate formatted expedition_code: "0" + seq if 1 digit, else just seq
CREATE OR REPLACE FUNCTION generate_expedition_code()
RETURNS VARCHAR AS $$
DECLARE
  seq_val BIGINT;
  formatted_code VARCHAR;
BEGIN
  seq_val := nextval('expedition_code_seq');
  IF seq_val < 10 THEN
    formatted_code := '0' || seq_val::VARCHAR;
  ELSE
    formatted_code := seq_val::VARCHAR;
  END IF;
  RETURN formatted_code;
END;
$$ LANGUAGE plpgsql;

-- Create table expeditions
CREATE TABLE IF NOT EXISTS expeditions (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  expedition_code VARCHAR(255) NOT NULL UNIQUE DEFAULT generate_expedition_code(),
  expedition_name VARCHAR(255),
  address VARCHAR(255),
  telp_number VARCHAR(255),
  phone_number VARCHAR(255),
  notes TEXT,
  created_at TIMESTAMP,
  created_by VARCHAR(255),
  updated_at TIMESTAMP,
  updated_by VARCHAR(255),
  deleted_at TIMESTAMP,
  deleted_by VARCHAR(255)
);

-- Indexes
CREATE INDEX IF NOT EXISTS expeditions_id_index ON expeditions (id);
CREATE INDEX IF NOT EXISTS expeditions_expedition_code_index ON expeditions (expedition_code);
CREATE INDEX IF NOT EXISTS expeditions_expedition_name_index ON expeditions (expedition_name);
CREATE INDEX IF NOT EXISTS expeditions_created_at_index ON expeditions (created_at);
CREATE INDEX IF NOT EXISTS expeditions_updated_at_index ON expeditions (updated_at);
CREATE INDEX IF NOT EXISTS expeditions_deleted_at_index ON expeditions (deleted_at);

