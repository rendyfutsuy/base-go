-- Seed Permission Groups for Module "Supplier"
INSERT INTO "permission_groups" ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
-- Supplier
    ('f1a2b3c4-d5e6-4f7a-8b9c-0d1e2f3a4b5c', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'View', false, 'Have Full Access for View Supplier Sub-Module', 'Supplier'),
    ('a2b3c4d5-e6f7-4a8b-9c0d-1e2f3a4b5c6d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Create', false, 'Have Full Access for Create Supplier Sub-Module', 'Supplier'),
    ('b3c4d5e6-f7a8-4b9c-0d1e-2f3a4b5c6d7e', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update', false, 'Have Full Access for Update Supplier Sub-Module', 'Supplier'),
    ('c4d5e6f7-a8b9-4c0d-1e2f-3a4b5c6d7e8f', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete', false, 'Have Full Access for Delete Supplier Sub-Module', 'Supplier'),
    ('d5e6f7a8-b9c0-4d1e-2f3a-4b5c6d7e8f9a', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Export', false, 'Have Full Access for Export Supplier Sub-Module', 'Supplier')
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions for Module "Supplier"
INSERT INTO "permissions" (
    "id",
    "created_at",
    "updated_at",
    "name",
    "deletable"
)
VALUES
-- Supplier Permissions
    (
        'e6f7a8b9-c0d1-4e2f-3a4b-5c6d7e8f9a0b',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'supplier.view',
        false
    ),
    (
        'f7a8b9c0-d1e2-4f3a-4b5c-6d7e8f9a0b1c',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'supplier.create',
        false
    ),
    (
        'a8b9c0d1-e2f3-4a4b-5c6d-7e8f9a0b1c2d',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'supplier.update',
        false
    ),
    (
        'b9c0d1e2-f3a4-4b5c-6d7e-8f9a0b1c2d3e',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'supplier.delete',
        false
    ),
    (
        'c0d1e2f3-a4b5-4c6d-7e8f-9a0b1c2d3e4f',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'supplier.export',
        false
    )
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions Modules (Permission Groups <-> Permissions) for Module "Supplier"
INSERT INTO "permissions_modules" (
    "permission_group_id",
    "permission_id"
)
VALUES
-- Supplier Permission Scope
    -- View permission group -> api.master-data.supplier.view
    (
        'f1a2b3c4-d5e6-4f7a-8b9c-0d1e2f3a4b5c',
        'e6f7a8b9-c0d1-4e2f-3a4b-5c6d7e8f9a0b'
    ),
    -- Create permission group -> api.master-data.supplier.create
    (
        'a2b3c4d5-e6f7-4a8b-9c0d-1e2f3a4b5c6d',
        'f7a8b9c0-d1e2-4f3a-4b5c-6d7e8f9a0b1c'
    ),
    -- Update permission group -> api.master-data.supplier.update
    (
        'b3c4d5e6-f7a8-4b9c-0d1e-2f3a4b5c6d7e',
        'a8b9c0d1-e2f3-4a4b-5c6d-7e8f9a0b1c2d'
    ),
    -- Delete permission group -> api.master-data.supplier.delete
    (
        'c4d5e6f7-a8b9-4c0d-1e2f-3a4b5c6d7e8f',
        'b9c0d1e2-f3a4-4b5c-6d7e-8f9a0b1c2d3e'
    ),
    -- Export permission group -> api.master-data.supplier.export
    (
        'd5e6f7a8-b9c0-4d1e-2f3a-4b5c6d7e8f9a',
        'c0d1e2f3-a4b5-4c6d-7e8f-9a0b1c2d3e4f'
    )
ON CONFLICT DO NOTHING;

-- Assign Permission Groups to Super Admin Role
INSERT INTO "modules_roles" (
    "permission_group_id",
    "role_id"
)
VALUES
-- Supplier Module to Super Admin Role Scope BEGIN
    (   
        'f1a2b3c4-d5e6-4f7a-8b9c-0d1e2f3a4b5c',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        'a2b3c4d5-e6f7-4a8b-9c0d-1e2f3a4b5c6d',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        'b3c4d5e6-f7a8-4b9c-0d1e-2f3a4b5c6d7e',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        'c4d5e6f7-a8b9-4c0d-1e2f-3a4b5c6d7e8f',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        'd5e6f7a8-b9c0-4d1e-2f3a-4b5c6d7e8f9a',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    )
ON CONFLICT DO NOTHING;
-- Supplier Module to Super Admin Role Scope END

