-- Create the uuid-ossp extension if it does not exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS modules_roles(
  permission_group_id UUID,
  role_id UUID
);