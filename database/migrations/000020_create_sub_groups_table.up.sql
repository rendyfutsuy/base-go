-- Create sequence for subgroup_code starting from 1
CREATE SEQUENCE IF NOT EXISTS subgroup_code_seq START WITH 1 INCREMENT BY 1;

-- Create table sub_groups
CREATE TABLE IF NOT EXISTS sub_groups (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  goods_group_id UUID NOT NULL REFERENCES goods_group(id) ON DELETE RESTRICT,
  subgroup_code BIGINT NOT NULL UNIQUE DEFAULT nextval('subgroup_code_seq'),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  created_by VARCHAR(255),
  updated_at TIMESTAMP,
  updated_by VARCHAR(255),
  deleted_at TIMESTAMP,
  deleted_by VARCHAR(255)
);

-- Indexes
CREATE INDEX IF NOT EXISTS sub_groups_id_index ON sub_groups (id);
CREATE INDEX IF NOT EXISTS sub_groups_goods_group_id_index ON sub_groups (goods_group_id);
CREATE INDEX IF NOT EXISTS sub_groups_subgroup_code_index ON sub_groups (subgroup_code);
CREATE INDEX IF NOT EXISTS sub_groups_name_index ON sub_groups (name);
CREATE INDEX IF NOT EXISTS sub_groups_created_at_index ON sub_groups (created_at);
CREATE INDEX IF NOT EXISTS sub_groups_updated_at_index ON sub_groups (updated_at);
CREATE INDEX IF NOT EXISTS sub_groups_deleted_at_index ON sub_groups (deleted_at);

-- Unique constraint: name must be unique within a goods_group (excluding soft-deleted records)
CREATE UNIQUE INDEX IF NOT EXISTS subgroup_in_group ON sub_groups (goods_group_id, name) WHERE deleted_at IS NULL;

