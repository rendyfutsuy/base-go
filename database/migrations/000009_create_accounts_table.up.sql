CREATE TABLE IF NOT EXISTS accounts(
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

CREATE INDEX accounts_id_index ON accounts (id);

CREATE INDEX accounts_created_at_index ON accounts (created_at);

CREATE INDEX accounts_updated_at_index ON accounts (updated_at);

CREATE INDEX accounts_deleted_at_index ON accounts (deleted_at);

CREATE INDEX accounts_code_index ON accounts (code);