CREATE TABLE IF NOT EXISTS carriages(
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

CREATE INDEX carriages_id_index ON carriages (id);

CREATE INDEX carriages_created_at_index ON carriages (created_at);

CREATE INDEX carriages_updated_at_index ON carriages (updated_at);

CREATE INDEX carriages_deleted_at_index ON carriages (deleted_at);

CREATE INDEX carriages_code_index ON carriages (code);