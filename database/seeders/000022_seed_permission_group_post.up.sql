-- Seed Permission Groups for Module "Post"
INSERT INTO "permission_groups" 
    ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
    ('9f3c2a7e-6b41-4d8c-9a2f-1c7e5b8d2f10', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Create', false, 'Have Full Access for Create Course Sub-Module', 'Post'),
    ('2c6e9b14-3f8a-4a7d-b5c2-8d1e4f7a9b33', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update', false, 'Have Full Access for Update Course Sub-Module', 'Post'),
    ('7a1d5e3c-8b92-4f6a-91c7-3e2b4d8f6a21', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete', false, 'Have Full Access for Delete Course Sub-Module', 'Post')
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions for Module "Post"
INSERT INTO "permissions" 
    ("id", "created_at", "updated_at", "name", "deletable")
VALUES
    ('4b8e2c71-9d53-4f1a-a6c8-2e7d5b3f9a44', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'post.create', false),
    ('8c2f6a19-1e4b-4d7c-b3a8-5f9d2e6c7a55', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'post.update', false),
    ('1e7b3c5a-6d9f-4a2e-8c1b-7d3f5a9c2e66', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'post.delete', false)
ON CONFLICT (id) DO NOTHING;

-- Map Permission Groups <-> Permissions
INSERT INTO "permissions_modules" 
    ("permission_group_id", "permission_id")
VALUES
    ('9f3c2a7e-6b41-4d8c-9a2f-1c7e5b8d2f10', '4b8e2c71-9d53-4f1a-a6c8-2e7d5b3f9a44'),
    ('2c6e9b14-3f8a-4a7d-b5c2-8d1e4f7a9b33', '8c2f6a19-1e4b-4d7c-b3a8-5f9d2e6c7a55'),
    ('7a1d5e3c-8b92-4f6a-91c7-3e2b4d8f6a21', '1e7b3c5a-6d9f-4a2e-8c1b-7d3f5a9c2e66')
ON CONFLICT DO NOTHING;

-- Assign Permission Groups to Super Admin Role
INSERT INTO "modules_roles" (
    "permission_group_id",
    "role_id"
)
VALUES
    ('9f3c2a7e-6b41-4d8c-9a2f-1c7e5b8d2f10', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('2c6e9b14-3f8a-4a7d-b5c2-8d1e4f7a9b33', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('7a1d5e3c-8b92-4f6a-91c7-3e2b4d8f6a21', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0')
ON CONFLICT DO NOTHING;