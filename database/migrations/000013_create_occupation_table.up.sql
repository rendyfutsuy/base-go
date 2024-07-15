CREATE TABLE IF NOT EXISTS occupations(
    id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
    subcob_id UUID NOT NULL,
    code VARCHAR(80) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP,
    created_by VARCHAR(80),
    updated_at TIMESTAMP,
    updated_by VARCHAR(80),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(80),

    CONSTRAINT fk_subcob FOREIGN KEY (subcob_id) REFERENCES subcobs(id) ON DELETE RESTRICT
);

CREATE INDEX occupations_id_index ON occupations (id);

CREATE INDEX occupations_created_at_index ON occupations (created_at);

CREATE INDEX occupations_updated_at_index ON occupations (updated_at);

CREATE INDEX occupations_deleted_at_index ON occupations (deleted_at);

CREATE INDEX occupations_code_index ON occupations (code);