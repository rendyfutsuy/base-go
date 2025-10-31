-- Create the uuid-ossp extension if it does not exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS permissions_modules(
  permission_id UUID,
  permission_group_id UUID
);