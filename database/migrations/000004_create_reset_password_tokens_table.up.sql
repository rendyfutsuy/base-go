-- Create the uuid-ossp extension if it does not exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS reset_password_tokens (
   created_at TIMESTAMP,
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP,
   user_id UUID NOT NULL, 
   access_token VARCHAR(225) NOT NULL,
   CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX reset_password_tokens_created_at_index ON reset_password_tokens (created_at);
CREATE INDEX reset_password_tokens_updated_at_index ON reset_password_tokens (updated_at);
CREATE INDEX reset_password_tokens_deleted_at_index ON reset_password_tokens (deleted_at);
CREATE INDEX reset_password_tokens_user_id_index ON reset_password_tokens (user_id);
