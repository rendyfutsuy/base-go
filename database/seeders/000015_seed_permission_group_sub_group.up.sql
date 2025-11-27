-- Seed Permission Groups for Module "Sub Golongan"
INSERT INTO "permission_groups" ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
    ('08bb3762-44a6-4908-8c10-b60b0312ddbe', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'View', false, 'Have Full Access for View Sub Golongan Sub-Module', 'Sub Golongan'),
    ('0ab68d1e-f41e-4171-bf85-9cc444af434a', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Create', false, 'Have Full Access for Create Sub Golongan Sub-Module', 'Sub Golongan'),
    ('eeccf1be-525e-460a-813e-ca45f9939477', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update', false, 'Have Full Access for Update Sub Golongan Sub-Module', 'Sub Golongan'),
    ('021a0b51-7738-44d7-98d6-44792d06c1e6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete', false, 'Have Full Access for Delete Sub Golongan Sub-Module', 'Sub Golongan'),
    ('64e19127-09f9-4639-aa05-faaca5156f8e', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Export', false, 'Have Full Access for Export Sub Golongan Sub-Module', 'Sub Golongan')
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions for Module "Sub Golongan"
INSERT INTO "permissions" (
    "id",
    "created_at",
    "updated_at",
    "name",
    "deletable"
)
VALUES
    ('36bf02f8-d507-4130-8bbf-e66f460de6a2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'sub-group.view', false),
    ('f366d0d5-2baa-43f6-9fc8-56afa136977d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'sub-group.create', false),
    ('4a519b42-d2a1-4689-b50f-5eb827324fe2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'sub-group.update', false),
    ('ce85a799-641c-4a8f-96ea-0bd6bb1561f7', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'sub-group.delete', false),
    ('3e0047fd-ac1f-4c91-b029-5e59e546eede', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'sub-group.export', false)
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions Modules (Permission Groups <-> Permissions) for Module "Sub Golongan"
INSERT INTO "permissions_modules" (
    "permission_group_id",
    "permission_id"
)
VALUES
    ('08bb3762-44a6-4908-8c10-b60b0312ddbe', '36bf02f8-d507-4130-8bbf-e66f460de6a2'),
    ('0ab68d1e-f41e-4171-bf85-9cc444af434a', 'f366d0d5-2baa-43f6-9fc8-56afa136977d'),
    ('eeccf1be-525e-460a-813e-ca45f9939477', '4a519b42-d2a1-4689-b50f-5eb827324fe2'),
    ('021a0b51-7738-44d7-98d6-44792d06c1e6', 'ce85a799-641c-4a8f-96ea-0bd6bb1561f7'),
    ('64e19127-09f9-4639-aa05-faaca5156f8e', '3e0047fd-ac1f-4c91-b029-5e59e546eede')
ON CONFLICT DO NOTHING;

-- Assign Permission Groups to Super Admin Role
INSERT INTO "modules_roles" (
    "permission_group_id",
    "role_id"
)
VALUES
-- Sub Golongan Module to Super Admin Role Scope BEGIN
    ('08bb3762-44a6-4908-8c10-b60b0312ddbe', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('0ab68d1e-f41e-4171-bf85-9cc444af434a', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('eeccf1be-525e-460a-813e-ca45f9939477', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('021a0b51-7738-44d7-98d6-44792d06c1e6', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('64e19127-09f9-4639-aa05-faaca5156f8e', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0')
ON CONFLICT DO NOTHING;
-- Sub Golongan Module to Super Admin Role Scope END

