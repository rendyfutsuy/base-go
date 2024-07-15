CREATE TABLE IF NOT EXISTS contractors(
    id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
    name VARCHAR(80) NOT NULL,
    code VARCHAR(80) NOT NULL UNIQUE,
    address TEXT NOT NULL,
    created_at TIMESTAMP,
    created_by VARCHAR(80),
    updated_at TIMESTAMP,
    updated_by VARCHAR(80),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(80)
);

CREATE INDEX contractors_id_index ON contractors (id);

CREATE INDEX contractors_created_at_index ON contractors (created_at);

CREATE INDEX contractors_updated_at_index ON contractors (updated_at);

CREATE INDEX contractors_deleted_at_index ON contractors (deleted_at);

CREATE INDEX contractors_code_index ON contractors (code);