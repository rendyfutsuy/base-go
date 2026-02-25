CREATE TABLE IF NOT EXISTS files (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  name VARCHAR(255) NOT NULL,
  file_path TEXT,
  description TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS file_id_index ON files (id);
CREATE INDEX IF NOT EXISTS file_name_index ON files (name);
CREATE INDEX IF NOT EXISTS file_created_at_index ON files (created_at);
CREATE INDEX IF NOT EXISTS file_updated_at_index ON files (updated_at);
CREATE INDEX IF NOT EXISTS file_deleted_at_index ON files (deleted_at);
