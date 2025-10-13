-- Create the uuid-ossp extension if it does not exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS permissions(
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP,
  name VARCHAR(255) NOT NULL UNIQUE,
  deletable boolean NOT NULL
);

CREATE INDEX permissions_id_index ON permissions (id);

CREATE INDEX permissions_created_at_index ON permissions (created_at);

CREATE INDEX permissions_updated_at_index ON permissions (updated_at);

CREATE INDEX permissions_deleted_at_index ON permissions (deleted_at);

CREATE INDEX permissions_name_index ON permissions (name);