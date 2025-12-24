-- Create sequence for backing_code starting from 1
CREATE SEQUENCE IF NOT EXISTS backing_code_seq START WITH 1 INCREMENT BY 1;

-- Function to generate formatted backing_code: "0" + seq if 1 digit, else just seq
CREATE OR REPLACE FUNCTION generate_backing_code()
RETURNS VARCHAR AS $$
DECLARE
  seq_val BIGINT;
  formatted_code VARCHAR;
BEGIN
  seq_val := nextval('backing_code_seq');
  IF seq_val < 10 THEN
    formatted_code := '0' || seq_val::VARCHAR;
  ELSE
    formatted_code := seq_val::VARCHAR;
  END IF;
  RETURN formatted_code;
END;
$$ LANGUAGE plpgsql;

-- Create table backings
CREATE TABLE IF NOT EXISTS backings (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  type_id UUID NOT NULL REFERENCES types(id),
  backing_code VARCHAR(255) NOT NULL UNIQUE DEFAULT generate_backing_code(),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  created_by VARCHAR(255),
  updated_at TIMESTAMP,
  updated_by VARCHAR(255),
  deleted_at TIMESTAMP,
  deleted_by VARCHAR(255)
);

-- Create unique index on (type_id, name) to ensure unique name within type
CREATE UNIQUE INDEX IF NOT EXISTS backing_in_type ON backings (type_id, name) WHERE deleted_at IS NULL;

-- Indexes
CREATE INDEX IF NOT EXISTS backings_id_index ON backings (id);
CREATE INDEX IF NOT EXISTS backings_type_id_index ON backings (type_id);
CREATE INDEX IF NOT EXISTS backings_backing_code_index ON backings (backing_code);
CREATE INDEX IF NOT EXISTS backings_name_index ON backings (name);
CREATE INDEX IF NOT EXISTS backings_created_at_index ON backings (created_at);
CREATE INDEX IF NOT EXISTS backings_updated_at_index ON backings (updated_at);
CREATE INDEX IF NOT EXISTS backings_deleted_at_index ON backings (deleted_at);
