CREATE TABLE IF NOT EXISTS subcobs(
    id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
    category_id UUID NOT NULL,
    cob_id UUID NOT NULL,
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

    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,

    CONSTRAINT fk_cob FOREIGN KEY (cob_id) REFERENCES cobs(id) ON DELETE RESTRICT
);

CREATE INDEX subcobs_id_index ON subcobs (id);

CREATE INDEX subcobs_created_at_index ON subcobs (created_at);

CREATE INDEX subcobs_updated_at_index ON subcobs (updated_at);

CREATE INDEX subcobs_deleted_at_index ON subcobs (deleted_at);

CREATE INDEX subcobs_code_index ON subcobs (code);