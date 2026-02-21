-- Create the uuid-ossp extension if it does not exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS otps (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   created_at TIMESTAMP,
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP,
   user_id UUID NOT NULL,
   token VARCHAR(255) NOT NULL,
   CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS otps_created_at_index ON otps (created_at);
CREATE INDEX IF NOT EXISTS otps_updated_at_index ON otps (updated_at);
CREATE INDEX IF NOT EXISTS otps_deleted_at_index ON otps (deleted_at);
CREATE INDEX IF NOT EXISTS otps_user_id_index ON otps (user_id);
CREATE INDEX IF NOT EXISTS otps_token_index ON otps (token);
