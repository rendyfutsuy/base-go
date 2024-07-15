-- Create the uuid-ossp extension if it does not exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS roles(
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  created_at TIMESTAMP,
  created_by VARCHAR(80),
  updated_at TIMESTAMP,
  updated_by VARCHAR(80),
  deleted_at TIMESTAMP,
  deleted_by VARCHAR(80),
  name VARCHAR(80) NOT NULL UNIQUE,
  deletable boolean NOT NULL
);

CREATE INDEX roles_id_index ON roles (id);

CREATE INDEX roles_created_at_index ON roles (created_at);

CREATE INDEX roles_updated_at_index ON roles (updated_at);

CREATE INDEX roles_deleted_at_index ON roles (deleted_at);

CREATE INDEX roles_name_index ON roles (name);