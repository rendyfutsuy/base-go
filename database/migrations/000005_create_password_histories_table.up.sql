-- Create the uuid-ossp extension if it does not exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS password_histories (
   created_at TIMESTAMP,
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP,
   user_id UUID NOT NULL, 
   hashed_password VARCHAR(225) NOT NULL,
   CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX password_histories_created_at_index ON password_histories (created_at);
CREATE INDEX password_histories_updated_at_index ON password_histories (updated_at);
CREATE INDEX password_histories_deleted_at_index ON password_histories (deleted_at);
CREATE INDEX password_histories_user_id_index ON password_histories (user_id);
