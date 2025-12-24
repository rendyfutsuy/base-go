-- Create sequence for type_code starting from 1
CREATE SEQUENCE IF NOT EXISTS type_code_seq START WITH 1 INCREMENT BY 1;

-- Function to generate type_code from sequence
CREATE OR REPLACE FUNCTION generate_type_code()
RETURNS VARCHAR AS $$
DECLARE
  seq_val BIGINT;
  formatted_code VARCHAR;
BEGIN
  seq_val := nextval('type_code_seq');
  IF seq_val < 10 THEN
    formatted_code := '0' || seq_val::VARCHAR;
  ELSE
    formatted_code := seq_val::VARCHAR;
  END IF;
  RETURN formatted_code;
END;
$$ LANGUAGE plpgsql;

-- Create table types
CREATE TABLE IF NOT EXISTS types (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  subgroup_id UUID NOT NULL REFERENCES sub_groups(id),
  type_code VARCHAR(255) NOT NULL UNIQUE DEFAULT generate_type_code(),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(255),
  updated_at TIMESTAMP NOT NULL,
  updated_by VARCHAR(255),
  deleted_at TIMESTAMP,
  deleted_by VARCHAR(255),
  CONSTRAINT type_in_subgroup UNIQUE (subgroup_id, name)
);

-- Indexes
CREATE INDEX IF NOT EXISTS types_id_index ON types (id);
CREATE INDEX IF NOT EXISTS types_subgroup_id_index ON types (subgroup_id);
CREATE INDEX IF NOT EXISTS types_type_code_index ON types (type_code);
CREATE INDEX IF NOT EXISTS types_name_index ON types (name);
CREATE INDEX IF NOT EXISTS types_created_at_index ON types (created_at);
CREATE INDEX IF NOT EXISTS types_updated_at_index ON types (updated_at);
CREATE INDEX IF NOT EXISTS types_deleted_at_index ON types (deleted_at);
