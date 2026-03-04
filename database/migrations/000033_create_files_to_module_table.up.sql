CREATE TABLE IF NOT EXISTS files_to_module (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  file_id UUID REFERENCES files(id),
  module_type VARCHAR(255) NOT NULL,
  module_id UUID NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS files_to_module_id_index ON files_to_module (id);
CREATE INDEX IF NOT EXISTS files_to_module_file_id_index ON files_to_module (file_id);
CREATE INDEX IF NOT EXISTS files_to_module_module_type_index ON files_to_module (module_type);
CREATE INDEX IF NOT EXISTS files_to_module_module_id_index ON files_to_module (module_id);
CREATE INDEX IF NOT EXISTS files_to_module_created_at_index ON files_to_module (created_at);
CREATE INDEX IF NOT EXISTS files_to_module_updated_at_index ON files_to_module (updated_at);
