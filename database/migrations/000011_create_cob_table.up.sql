CREATE TABLE IF NOT EXISTS cobs(
    id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
    category_id UUID NOT NULL,
    name VARCHAR(80) NOT NULL,
    code VARCHAR(80) NOT NULL,
    forms TEXT,
    active_date TIMESTAMP,
    is_hidden_from_facultative BOOLEAN,
    is_inactive BOOLEAN,
    is_from_web_credit BOOLEAN,
    created_at TIMESTAMP,
    created_by VARCHAR(80),
    updated_at TIMESTAMP,
    updated_by VARCHAR(80),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(80),

    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);

CREATE INDEX cobs_id_index ON cobs (id);

CREATE INDEX cobs_created_at_index ON cobs (created_at);

CREATE INDEX cobs_updated_at_index ON cobs (updated_at);

CREATE INDEX cobs_deleted_at_index ON cobs (deleted_at);

CREATE INDEX cobs_code_index ON cobs (code);