CREATE TABLE IF NOT EXISTS conveyances(
    id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
    name VARCHAR(80) NOT NULL,
    code VARCHAR(80) NOT NULL UNIQUE,
    type VARCHAR(80) NOT NULL,
    created_at TIMESTAMP,
    created_by VARCHAR(80),
    updated_at TIMESTAMP,
    updated_by VARCHAR(80),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(80)
);

CREATE INDEX conveyances_id_index ON conveyances (id);

CREATE INDEX conveyances_created_at_index ON conveyances (created_at);

CREATE INDEX conveyances_updated_at_index ON conveyances (updated_at);

CREATE INDEX conveyances_deleted_at_index ON conveyances (deleted_at);

CREATE INDEX conveyances_code_index ON conveyances (code);