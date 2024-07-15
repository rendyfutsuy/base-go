CREATE TABLE IF NOT EXISTS categories(
    id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
    name VARCHAR(80) NOT NULL,
    code VARCHAR(80) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP,
    created_by VARCHAR(80),
    updated_at TIMESTAMP,
    updated_by VARCHAR(80),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(80)
);

CREATE INDEX categories_id_index ON categories (id);

CREATE INDEX categories_created_at_index ON categories (created_at);

CREATE INDEX categories_updated_at_index ON categories (updated_at);

CREATE INDEX categories_deleted_at_index ON categories (deleted_at);

CREATE INDEX categories_code_index ON categories (code);