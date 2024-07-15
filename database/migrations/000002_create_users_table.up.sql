-- Create the uuid-ossp extension if it does not exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users(
   id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
   created_at TIMESTAMP,
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP,
   role_id UUID NOT NULL, 
   nik VARCHAR(80) NOT NULL UNIQUE,
   username VARCHAR(80) NOT NULL,
   full_name VARCHAR(80) NOT NULL,
   email VARCHAR(80) NOT NULL UNIQUE,
   password VARCHAR(225) NOT NULL,
   -- new column BEGIN --
   is_active BOOLEAN NOT NULL DEFAULT TRUE,
   password_expired_at TIMESTAMP,
   counter INT NOT NULL DEFAULT 0,
   gender VARCHAR(15),
   -- new column END --
   CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);

CREATE INDEX users_id_index ON users (id);

CREATE INDEX users_created_at_index ON users (created_at);

CREATE INDEX users_updated_at_index ON users (updated_at);

CREATE INDEX users_deleted_at_index ON users (deleted_at);

CREATE INDEX users_role_id_index ON users (role_id);

CREATE INDEX users_full_name_index ON users (full_name);