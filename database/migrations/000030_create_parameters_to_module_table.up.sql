CREATE TABLE IF NOT EXISTS parameters_to_module (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  parameter_id UUID REFERENCES parameters(id),
  module_type VARCHAR(255) NOT NULL,
  module_id UUID NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS parameters_to_module_id_index ON parameters_to_module (id);
CREATE INDEX IF NOT EXISTS parameters_to_module_parameter_id_index ON parameters_to_module (parameter_id);
CREATE INDEX IF NOT EXISTS parameters_to_module_module_type_index ON parameters_to_module (module_type);
CREATE INDEX IF NOT EXISTS parameters_to_module_module_id_index ON parameters_to_module (module_id);
CREATE INDEX IF NOT EXISTS parameters_to_module_created_at_index ON parameters_to_module (created_at);
CREATE INDEX IF NOT EXISTS parameters_to_module_updated_at_index ON parameters_to_module (updated_at);
