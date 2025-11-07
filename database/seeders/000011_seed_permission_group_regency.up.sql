-- Seed Permission Groups for Module "Regency"
INSERT INTO "permission_groups" ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
    ('1a2b3c4d-e5f6-4a7b-8c9d-0e1f2a3b4c5d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'View', false, 'Have Full Access for View Regency Module', 'Regency'),
    ('2b3c4d5e-f6a7-4b8c-9d0e-1f2a3b4c5d6e', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Create', false, 'Have Full Access for Create Regency Module', 'Regency'),
    ('3c4d5e6f-a7b8-4c9d-0e1f-2a3b4c5d6e7f', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update', false, 'Have Full Access for Update Regency Module', 'Regency'),
    ('4d5e6f7a-b8c9-4d0e-1f2a-3b4c5d6e7f8a', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete', false, 'Have Full Access for Delete Regency Module', 'Regency'),
    ('5e6f7a8b-c9d0-4e1f-2a3b-4c5d6e7f8a9b', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Export', false, 'Have Full Access for Export Regency Module', 'Regency')
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions for Module "Regency"
INSERT INTO "permissions" (
    "id",
    "created_at",
    "updated_at",
    "name",
    "deletable"
)
VALUES
    ('6f7a8b9c-d0e1-4f2a-3b4c-5d6e7f8a9b0c', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'province.view', false),
    ('7a8b9c0d-e1f2-4a3b-4c5d-6e7f8a9b0c1d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'province.create', false),
    ('8b9c0d1e-f2a3-4b4c-5d6e-7f8a9b0c1d2e', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'province.update', false),
    ('9c0d1e2f-a3b4-4c5d-6e7f-8a9b0c1d2e3f', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'province.delete', false),
    ('0d1e2f3a-b4c5-4d6e-7f8a-9b0c1d2e3f4a', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'province.export', false)
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions Modules (Permission Groups <-> Permissions) for Module "Regency"
INSERT INTO "permissions_modules" (
    "permission_group_id",
    "permission_id"
)
VALUES
    ('1a2b3c4d-e5f6-4a7b-8c9d-0e1f2a3b4c5d', '6f7a8b9c-d0e1-4f2a-3b4c-5d6e7f8a9b0c'),
    ('2b3c4d5e-f6a7-4b8c-9d0e-1f2a3b4c5d6e', '7a8b9c0d-e1f2-4a3b-4c5d-6e7f8a9b0c1d'),
    ('3c4d5e6f-a7b8-4c9d-0e1f-2a3b4c5d6e7f', '8b9c0d1e-f2a3-4b4c-5d6e-7f8a9b0c1d2e'),
    ('4d5e6f7a-b8c9-4d0e-1f2a-3b4c5d6e7f8a', '9c0d1e2f-a3b4-4c5d-6e7f-8a9b0c1d2e3f'),
    ('5e6f7a8b-c9d0-4e1f-2a3b-4c5d6e7f8a9b', '0d1e2f3a-b4c5-4d6e-7f8a-9b0c1d2e3f4a')
ON CONFLICT DO NOTHING;

-- Assign Permission Groups to Super Admin Role
INSERT INTO "modules_roles" (
    "permission_group_id",
    "role_id"
)
VALUES
-- Regency Module to Super Admin Role Scope BEGIN
    ('1a2b3c4d-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('2b3c4d5e-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('3c4d5e6f-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('4d5e6f7a-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('5e6f7a8b-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0')
ON CONFLICT DO NOTHING;
-- Regency Module to Super Admin Role Scope END

