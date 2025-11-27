-- Seed Permission Groups for Module "Jenis"
INSERT INTO "permission_groups" ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
    ('b8a10397-076d-4a16-9b6b-de4e84dcf420', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'View', false, 'Have Full Access for View Jenis Sub-Module', 'Jenis'),
    ('2ba54fa7-f54f-4110-bd2b-dbc9e0e16e43', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Create', false, 'Have Full Access for Create Jenis Sub-Module', 'Jenis'),
    ('05b985f4-f13f-4f98-9515-8c427536ca16', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update', false, 'Have Full Access for Update Jenis Sub-Module', 'Jenis'),
    ('149006b0-a0fa-429d-8327-c6cdc2b2963d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete', false, 'Have Full Access for Delete Jenis Sub-Module', 'Jenis'),
    ('38cd0018-b33b-4ce3-a1f9-d7fadcfb2bf4', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Export', false, 'Have Full Access for Export Jenis Sub-Module', 'Jenis')
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions for Module "Jenis"
INSERT INTO "permissions" (
    "id",
    "created_at",
    "updated_at",
    "name",
    "deletable"
)
VALUES
    ('83ab9222-8bea-4237-a194-e49de405a975', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'type.view', false),
    ('79887555-64cb-402d-a147-bc538931c14f', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'type.create', false),
    ('47240d3e-bb18-4b90-ade9-70ec397d7756', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'type.update', false),
    ('436cc1dd-24ae-4d8d-8818-ff118ec0b8ff', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'type.delete', false),
    ('9049031a-8df1-4c35-93fe-ebd6b41921cc', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'type.export', false)
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions Modules (Permission Groups <-> Permissions) for Module "Jenis"
INSERT INTO "permissions_modules" (
    "permission_group_id",
    "permission_id"
)
VALUES
    ('b8a10397-076d-4a16-9b6b-de4e84dcf420', '83ab9222-8bea-4237-a194-e49de405a975'),
    ('2ba54fa7-f54f-4110-bd2b-dbc9e0e16e43', '79887555-64cb-402d-a147-bc538931c14f'),
    ('05b985f4-f13f-4f98-9515-8c427536ca16', '47240d3e-bb18-4b90-ade9-70ec397d7756'),
    ('149006b0-a0fa-429d-8327-c6cdc2b2963d', '436cc1dd-24ae-4d8d-8818-ff118ec0b8ff'),
    ('38cd0018-b33b-4ce3-a1f9-d7fadcfb2bf4', '9049031a-8df1-4c35-93fe-ebd6b41921cc')
ON CONFLICT DO NOTHING;

-- Assign Permission Groups to Super Admin Role
INSERT INTO "modules_roles" (
    "permission_group_id",
    "role_id"
)
VALUES
-- Jenis Module to Super Admin Role Scope BEGIN
    ('b8a10397-076d-4a16-9b6b-de4e84dcf420', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('2ba54fa7-f54f-4110-bd2b-dbc9e0e16e43', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('05b985f4-f13f-4f98-9515-8c427536ca16', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('149006b0-a0fa-429d-8327-c6cdc2b2963d', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0'),
    ('38cd0018-b33b-4ce3-a1f9-d7fadcfb2bf4', 'a43a5e5f-a172-42d1-a70e-8834bf653eb0')
ON CONFLICT DO NOTHING;
-- Jenis Module to Super Admin Role Scope END

