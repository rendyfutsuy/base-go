CREATE TABLE IF NOT EXISTS shipyards(
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
  name VARCHAR(80) NOT NULL,
  code VARCHAR(80) NOT NULL UNIQUE,
  yard VARCHAR(80) NOT NULL,
  created_at TIMESTAMP,
  created_by VARCHAR(80),
  updated_at TIMESTAMP,
  updated_by VARCHAR(80),
  deleted_at TIMESTAMP,
  deleted_by VARCHAR(80)
);

CREATE SEQUENCE IF NOT EXISTS shipyards_seq;

CREATE OR REPLACE FUNCTION increment_code()
RETURNS TRIGGER AS $$
BEGIN
  NEW.code := TO_CHAR(NEXTVAL('shipyards_seq'), 'FM0000000');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER shipyards_trigger
BEFORE INSERT ON shipyards
FOR EACH ROW
EXECUTE PROCEDURE increment_code();

CREATE INDEX IF NOT EXISTS shipyards_id_index ON shipyards (id);

CREATE INDEX IF NOT EXISTS shipyards_created_at_index ON shipyards (created_at);

CREATE INDEX IF NOT EXISTS shipyards_updated_at_index ON shipyards (updated_at);

CREATE INDEX IF NOT EXISTS shipyards_deleted_at_index ON shipyards (deleted_at);

CREATE INDEX IF NOT EXISTS shipyards_name_index ON shipyards (name);