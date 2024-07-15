CREATE TABLE IF NOT EXISTS classes(
    id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
    name VARCHAR(80) NOT NULL,
    code VARCHAR(80) NOT NULL UNIQUE,
    created_at TIMESTAMP,
    created_by VARCHAR(80),
    updated_at TIMESTAMP,
    updated_by VARCHAR(80),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(80)
);

CREATE INDEX classes_id_index ON classes (id);

CREATE INDEX classes_created_at_index ON classes (created_at);

CREATE INDEX classes_updated_at_index ON classes (updated_at);

CREATE INDEX classes_deleted_at_index ON classes (deleted_at);

CREATE INDEX classes_code_index ON classes (code);