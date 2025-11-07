-- Seed Permission Groups for Module "Parameter"
INSERT INTO "permission_groups" ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
-- Parameter
    ('a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'View', false, 'Have Full Access for View Parameter Sub-Module', 'Parameter'),
    ('b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Create', false, 'Have Full Access for Create Parameter Sub-Module', 'Parameter'),
    ('c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update', false, 'Have Full Access for Update Parameter Sub-Module', 'Parameter'),
    ('d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete', false, 'Have Full Access for Delete Parameter Sub-Module', 'Parameter'),
    ('e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Export', false, 'Have Full Access for Export Parameter Sub-Module', 'Parameter')
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions for Module "Parameter"
INSERT INTO "permissions" (
    "id",
    "created_at",
    "updated_at",
    "name",
    "deletable"
)
VALUES
-- Parameter Permissions
    (
        '1f2a3b4c-5d6e-4f7a-8b9c-0d1e2f3a4b5c',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'parameter.view',
        false
    ),
    (
        '2a3b4c5d-6e7f-4a8b-9c0d-1e2f3a4b5c6d',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'parameter.create',
        false
    ),
    (
        '3b4c5d6e-7f8a-4b9c-0d1e-2f3a4b5c6d7e',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'parameter.update',
        false
    ),
    (
        '4c5d6e7f-8a9b-4c0d-1e2f-3a4b5c6d7e8f',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'parameter.delete',
        false
    ),
    (
        '5d6e7f8a-9b0c-4d1e-2f3a-4b5c6d7e8f9a',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'parameter.export',
        false
    )
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions Modules (Permission Groups <-> Permissions) for Module "Parameter"
INSERT INTO "permissions_modules" (
    "permission_group_id",
    "permission_id"
)
VALUES
-- Parameter Permission Scope
    -- View permission group -> api.master-data.parameter.view
    (
        'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d',
        '1f2a3b4c-5d6e-4f7a-8b9c-0d1e2f3a4b5c'
    ),
    -- Create permission group -> api.master-data.parameter.create
    (
        'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e',
        '2a3b4c5d-6e7f-4a8b-9c0d-1e2f3a4b5c6d'
    ),
    -- Update permission group -> api.master-data.parameter.update
    (
        'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f',
        '3b4c5d6e-7f8a-4b9c-0d1e-2f3a4b5c6d7e'
    ),
    -- Delete permission group -> api.master-data.parameter.delete
    (
        'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a',
        '4c5d6e7f-8a9b-4c0d-1e2f-3a4b5c6d7e8f'
    ),
    -- Export permission group -> api.master-data.parameter.export
    (
        'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b',
        '5d6e7f8a-9b0c-4d1e-2f3a-4b5c6d7e8f9a'
    )
ON CONFLICT DO NOTHING;

-- Assign Permission Groups to Super Admin Role
INSERT INTO "modules_roles" (
    "permission_group_id",
    "role_id"
)
VALUES
-- Parameter Module to Super Admin Role Scope BEGIN
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
        'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    )
ON CONFLICT DO NOTHING;
-- Parameter Module to Super Admin Role Scope END

