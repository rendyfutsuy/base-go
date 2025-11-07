-- Seed Permission Groups for Module "Golongan"
INSERT INTO "permission_groups" ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
-- Golongan
    ('a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'View', false, 'Have Full Access for View Golongan Sub-Module', 'Golongan'),
    ('b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Create', false, 'Have Full Access for Create Golongan Sub-Module', 'Golongan'),
    ('c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update', false, 'Have Full Access for Update Golongan Sub-Module', 'Golongan'),
    ('d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete', false, 'Have Full Access for Delete Golongan Sub-Module', 'Golongan'),
    ('e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Export', false, 'Have Full Access for Export Golongan Sub-Module', 'Golongan');

-- Seed Permissions for Module "Golongan"
INSERT INTO "permissions" (
    "id",
    "created_at",
    "updated_at",
    "name",
    "deletable"
)
VALUES
-- Golongan Permissions
    (
        'f1a2b3c4-d5e6-4f7a-8b9c-0d1e2f3a4b5c',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'group.view',
        false
    ),
    (
        'a2b3c4d5-e6f7-4a8b-9c0d-1e2f3a4b5c6d',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'group.create',
        false
    ),
    (
        'b3c4d5e6-f7a8-4b9c-0d1e-2f3a4b5c6d7e',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'group.update',
        false
    ),
    (
        'c4d5e6f7-a8b9-4c0d-1e2f-3a4b5c6d7e8f',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'group.delete',
        false
    ),
    (
        'd5e6f7a8-b9c0-4d1e-2f3a-4b5c6d7e8f9a',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'group.export',
        false
    );

-- Seed Permissions Modules (Permission Groups <-> Permissions) for Module "Golongan"
INSERT INTO "permissions_modules" (
    "permission_group_id",
    "permission_id"
)
VALUES
-- Golongan Permission Scope
    -- View permission group -> api.master-data.group.view
    (
        'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d',
        'f1a2b3c4-d5e6-4f7a-8b9c-0d1e2f3a4b5c'
    ),
    -- Create permission group -> api.master-data.group.create
    (
        'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e',
        'a2b3c4d5-e6f7-4a8b-9c0d-1e2f3a4b5c6d'
    ),
    -- Update permission group -> api.master-data.group.update
    (
        'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f',
        'b3c4d5e6-f7a8-4b9c-0d1e-2f3a4b5c6d7e'
    ),
    -- Delete permission group -> api.master-data.group.delete
    (
        'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a',
        'c4d5e6f7-a8b9-4c0d-1e2f-3a4b5c6d7e8f'
    ),
    -- Export permission group -> api.master-data.group.export
    (
        'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b',
        'd5e6f7a8-b9c0-4d1e-2f3a-4b5c6d7e8f9a'
    );

--- update permission group to super admin role
INSERT INTO
    "modules_roles" (
        "permission_group_id",
        "role_id"
    )
VALUES
-- Golongan Module to Super Admin Role Scope BEGIN
    (   
        'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        'd2193e52-4e89-4b50-b7ab-745d6ab36a22',
        'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b'
    );
-- Golongan Module to Super Admin Role Scope END