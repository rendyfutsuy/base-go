-- Create table parameter
CREATE TABLE IF NOT EXISTS parameters (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  code VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  value VARCHAR(255),
  type VARCHAR(255),
  description TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS parameter_id_index ON parameters (id);
CREATE INDEX IF NOT EXISTS parameter_code_index ON parameters (code);
CREATE INDEX IF NOT EXISTS parameter_name_index ON parameters (name);
CREATE INDEX IF NOT EXISTS parameter_type_index ON parameters (type);
CREATE INDEX IF NOT EXISTS parameter_created_at_index ON parameters (created_at);
CREATE INDEX IF NOT EXISTS parameter_updated_at_index ON parameters (updated_at);
CREATE INDEX IF NOT EXISTS parameter_deleted_at_index ON parameters (deleted_at);

