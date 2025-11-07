-- Create table parameter
CREATE TABLE IF NOT EXISTS parameter (
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
CREATE INDEX IF NOT EXISTS parameter_id_index ON parameter (id);
CREATE INDEX IF NOT EXISTS parameter_code_index ON parameter (code);
CREATE INDEX IF NOT EXISTS parameter_name_index ON parameter (name);
CREATE INDEX IF NOT EXISTS parameter_type_index ON parameter (type);
CREATE INDEX IF NOT EXISTS parameter_created_at_index ON parameter (created_at);
CREATE INDEX IF NOT EXISTS parameter_updated_at_index ON parameter (updated_at);
CREATE INDEX IF NOT EXISTS parameter_deleted_at_index ON parameter (deleted_at);

