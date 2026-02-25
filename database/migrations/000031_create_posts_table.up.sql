CREATE TABLE IF NOT EXISTS courses (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  created_by UUID NOT NULL REFERENCES users(id),
  title VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  short_description VARCHAR(255) NOT NULL,
  price DECIMAL(18,2) NOT NULL,
  discount_rate DECIMAL(5,2) NOT NULL,
  thumbnail_url TEXT,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS courses_id_index ON courses (id);
CREATE INDEX IF NOT EXISTS courses_created_by_index ON courses (created_by);
CREATE INDEX IF NOT EXISTS courses_title_index ON courses (title);
CREATE INDEX IF NOT EXISTS courses_created_at_index ON courses (created_at);
CREATE INDEX IF NOT EXISTS courses_updated_at_index ON courses (updated_at);
