-- Seed Permission Group "Delete User" for Module "Users"
INSERT INTO "permission_groups" ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
    ('e6f7a8b9-c0d1-4e2f-3a4b-5c6d7e8f9a0b', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete User', false, 'Have Full Access for Delete User Sub-Module', 'Users')
ON CONFLICT (id) DO NOTHING;

-- Seed Permission "user.delete"
INSERT INTO "permissions" (
    "id",
    "created_at",
    "updated_at",
    "name",
    "deletable"
)
VALUES
    ('40bc8e07-bc88-40ad-ac7a-d323abaeeccb', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'user.delete', false)
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions Modules (Permission Groups <-> Permissions) for "Delete User" Permission Group
INSERT INTO "permissions_modules" (
    "permission_group_id",
    "permission_id"
)
VALUES
    ('e6f7a8b9-c0d1-4e2f-3a4b-5c6d7e8f9a0b', '40bc8e07-bc88-40ad-ac7a-d323abaeeccb')
ON CONFLICT DO NOTHING;

-- Assign Permission Group "Delete User" to Super Admin Role
INSERT INTO "modules_roles" (
    "permission_group_id",
    "role_id"
)
VALUES
    ('e6f7a8b9-c0d1-4e2f-3a4b-5c6d7e8f9a0b', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0')
ON CONFLICT DO NOTHING;

