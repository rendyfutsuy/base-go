-- Create sequence for group_code starting from 1
CREATE SEQUENCE IF NOT EXISTS groups_code_seq START WITH 1 INCREMENT BY 1;

-- Function to generate formatted group_code: "0" + seq if 1 digit, else just seq
CREATE OR REPLACE FUNCTION generate_group_code()
RETURNS VARCHAR AS $$
DECLARE
  seq_val BIGINT;
  formatted_code VARCHAR;
BEGIN
  seq_val := nextval('groups_code_seq');
  IF seq_val < 10 THEN
    formatted_code := '0' || seq_val::VARCHAR;
  ELSE
    formatted_code := seq_val::VARCHAR;
  END IF;
  RETURN formatted_code;
END;
$$ LANGUAGE plpgsql;

-- Create table groups
CREATE TABLE IF NOT EXISTS groups (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  group_code VARCHAR(255) NOT NULL UNIQUE DEFAULT generate_group_code(), -- auto generate: "01", "02", ..., "09", "10", "11", ...
  name VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS groups_id_index ON groups (id);
CREATE INDEX IF NOT EXISTS groups_created_at_index ON groups (created_at);
CREATE INDEX IF NOT EXISTS groups_updated_at_index ON groups (updated_at);
CREATE INDEX IF NOT EXISTS groups_deleted_at_index ON groups (deleted_at);
CREATE INDEX IF NOT EXISTS groups_name_index ON groups (name);
CREATE INDEX IF NOT EXISTS groups_group_code_index ON groups (group_code);

