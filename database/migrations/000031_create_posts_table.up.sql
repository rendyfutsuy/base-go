CREATE TABLE IF NOT EXISTS posts (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  created_by UUID NOT NULL REFERENCES users(id),
  title VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  short_description VARCHAR(255) NOT NULL,
  thumbnail_url TEXT,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS posts_id_index ON posts (id);
CREATE INDEX IF NOT EXISTS posts_created_by_index ON posts (created_by);
CREATE INDEX IF NOT EXISTS posts_title_index ON posts (title);
CREATE INDEX IF NOT EXISTS posts_created_at_index ON posts (created_at);
CREATE INDEX IF NOT EXISTS posts_updated_at_index ON posts (updated_at);
