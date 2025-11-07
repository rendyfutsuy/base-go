-- Seed Permission Groups for Module "Expedition"
INSERT INTO "permission_groups" ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
    ('e1f2a3b4-c5d6-4e7f-8a9b-0c1d2e3f4a5b', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'View', false, 'Have Full Access for View Expedition Sub-Module', 'Expedition'),
    ('f2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Create', false, 'Have Full Access for Create Expedition Sub-Module', 'Expedition'),
    ('a3b4c5d6-e7f8-4a9b-0c1d-2e3f4a5b6c7d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update', false, 'Have Full Access for Update Expedition Sub-Module', 'Expedition'),
    ('b4c5d6e7-f8a9-4b0c-1d2e-3f4a5b6c7d8e', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete', false, 'Have Full Access for Delete Expedition Sub-Module', 'Expedition'),
    ('c5d6e7f8-a9b0-4c1d-2e3f-4a5b6c7d8e9f', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Export', false, 'Have Full Access for Export Expedition Sub-Module', 'Expedition')
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions for Module "Expedition"
INSERT INTO "permissions" (
    "id",
    "created_at",
    "updated_at",
    "name",
    "deletable"
)
VALUES
    ('d6e7f8a9-b0c1-4d2e-3f4a-5b6c7d8e9f0a', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'expedition.view', false),
    ('e7f8a9b0-c1d2-4e3f-4a5b-6c7d8e9f0a1b', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'expedition.create', false),
    ('f8a9b0c1-d2e3-4f4a-5b6c-7d8e9f0a1b2c', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'expedition.update', false),
    ('a9b0c1d2-e3f4-4a5b-6c7d-8e9f0a1b2c3d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'expedition.delete', false),
    ('b0c1d2e3-f4a5-4b6c-7d8e-9f0a1b2c3d4e', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'expedition.export', false)
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions Modules (Permission Groups <-> Permissions) for Module "Expedition"
INSERT INTO "permissions_modules" (
    "permission_group_id",
    "permission_id"
)
VALUES
    ('e1f2a3b4-c5d6-4e7f-8a9b-0c1d2e3f4a5b', 'd6e7f8a9-b0c1-4d2e-3f4a-5b6c7d8e9f0a'),
    ('f2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c', 'e7f8a9b0-c1d2-4e3f-4a5b-6c7d8e9f0a1b'),
    ('a3b4c5d6-e7f8-4a9b-0c1d-2e3f4a5b6c7d', 'f8a9b0c1-d2e3-4f4a-5b6c-7d8e9f0a1b2c'),
    ('b4c5d6e7-f8a9-4b0c-1d2e-3f4a5b6c7d8e', 'a9b0c1d2-e3f4-4a5b-6c7d-8e9f0a1b2c3d'),
    ('c5d6e7f8-a9b0-4c1d-2e3f-4a5b6c7d8e9f', 'b0c1d2e3-f4a5-4b6c-7d8e-9f0a1b2c3d4e')
ON CONFLICT DO NOTHING;

-- Assign Permission Groups to Super Admin Role
INSERT INTO "modules_roles" (
    "permission_group_id",
    "role_id"
)
VALUES
-- Expedition Module to Super Admin Role Scope BEGIN
    ('e1f2a3b4-c5d6-4e7f-8a9b-0c1d2e3f4a5b', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('f2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('a3b4c5d6-e7f8-4a9b-0c1d-2e3f4a5b6c7d', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('b4c5d6e7-f8a9-4b0c-1d2e-3f4a5b6c7d8e', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('c5d6e7f8-a9b0-4c1d-2e3f-4a5b6c7d8e9f', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0')
ON CONFLICT DO NOTHING;
-- Expedition Module to Super Admin Role Scope END

