-- Create sequence for group_code starting from 1
CREATE SEQUENCE IF NOT EXISTS goods_group_code_seq START WITH 1 INCREMENT BY 1;

-- Function to generate formatted group_code: "0" + seq if 1 digit, else just seq
CREATE OR REPLACE FUNCTION generate_group_code()
RETURNS VARCHAR AS $$
DECLARE
  seq_val BIGINT;
  formatted_code VARCHAR;
BEGIN
  seq_val := nextval('goods_group_code_seq');
  IF seq_val < 10 THEN
    formatted_code := '0' || seq_val::VARCHAR;
  ELSE
    formatted_code := seq_val::VARCHAR;
  END IF;
  RETURN formatted_code;
END;
$$ LANGUAGE plpgsql;

-- Create table goods_group
CREATE TABLE IF NOT EXISTS goods_group (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  group_code VARCHAR(255) NOT NULL UNIQUE DEFAULT generate_group_code(), -- auto generate: "01", "02", ..., "09", "10", "11", ...
  name VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS goods_group_id_index ON goods_group (id);
CREATE INDEX IF NOT EXISTS goods_group_created_at_index ON goods_group (created_at);
CREATE INDEX IF NOT EXISTS goods_group_updated_at_index ON goods_group (updated_at);
CREATE INDEX IF NOT EXISTS goods_group_deleted_at_index ON goods_group (deleted_at);
CREATE INDEX IF NOT EXISTS goods_group_name_index ON goods_group (name);
CREATE INDEX IF NOT EXISTS goods_group_group_code_index ON goods_group (group_code);

