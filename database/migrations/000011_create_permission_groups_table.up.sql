-- Create the uuid-ossp extension if it does not exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS permission_groups(
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP,
  name VARCHAR(255) NOT NULL,
  module VARCHAR(255) NOT NULL,
  description Text,
  deletable boolean NOT NULL
);

CREATE INDEX permission_groups_id_index ON permission_groups (id);

CREATE INDEX permission_groups_created_at_index ON permission_groups (created_at);

CREATE INDEX permission_groups_updated_at_index ON permission_groups (updated_at);

CREATE INDEX permission_groups_deleted_at_index ON permission_groups (deleted_at);

CREATE INDEX permission_groups_name_index ON permission_groups (name);