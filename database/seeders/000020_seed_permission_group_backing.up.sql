-- Seed Permission Groups for Module "Backing"
INSERT INTO "permission_groups" ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
-- Backing
    ('51dab038-4d4c-4d7f-8b94-2f3d642d13ed', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'View', false, 'Have Full Access for View Backing Sub-Module', 'Backing'),
    ('557777a3-022d-4013-985e-26c60a7e589a', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Create', false, 'Have Full Access for Create Backing Sub-Module', 'Backing'),
    ('e13d1dbf-d235-45ff-8e54-b0ad2c8fbaa3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update', false, 'Have Full Access for Update Backing Sub-Module', 'Backing'),
    ('59a0721f-08cb-4e12-a23b-bdbbab47abdc', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete', false, 'Have Full Access for Delete Backing Sub-Module', 'Backing'),
    ('39cc8499-e573-404f-a301-f0806988b0e9', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Export', false, 'Have Full Access for Export Backing Sub-Module', 'Backing')
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions for Module "Backing"
INSERT INTO "permissions" (
    "id",
    "created_at",
    "updated_at",
    "name",
    "deletable"
)
VALUES
-- Backing Permissions
    (
        '2ba6c0e3-c935-4856-bcac-4fe56b8d975d',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'backing.view',
        false
    ),
    (
        '1a27a520-b66d-4b84-b67c-18b009cb8754',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'backing.create',
        false
    ),
    (
        '9c5b5801-7602-4631-bbee-86a878d414c9',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'backing.update',
        false
    ),
    (
        'c69a9c7e-aac9-4463-8a9a-04f1ecc55a02',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'backing.delete',
        false
    ),
    (
        'ed2c4b27-f1a7-4552-86ae-0f0c7517bbda',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        'backing.export',
        false
    )
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions Modules (Permission Groups <-> Permissions) for Module "Backing"
INSERT INTO "permissions_modules" (
    "permission_group_id",
    "permission_id"
)
VALUES
-- Backing Permission Scope
    -- View permission group -> backing.view
    (
        '51dab038-4d4c-4d7f-8b94-2f3d642d13ed',
        '2ba6c0e3-c935-4856-bcac-4fe56b8d975d'
    ),
    -- Create permission group -> backing.create
    (
        '557777a3-022d-4013-985e-26c60a7e589a',
        '1a27a520-b66d-4b84-b67c-18b009cb8754'
    ),
    -- Update permission group -> backing.update
    (
        'e13d1dbf-d235-45ff-8e54-b0ad2c8fbaa3',
        '9c5b5801-7602-4631-bbee-86a878d414c9'
    ),
    -- Delete permission group -> backing.delete
    (
        '59a0721f-08cb-4e12-a23b-bdbbab47abdc',
        'c69a9c7e-aac9-4463-8a9a-04f1ecc55a02'
    ),
    -- Export permission group -> backing.export
    (
        '39cc8499-e573-404f-a301-f0806988b0e9',
        'ed2c4b27-f1a7-4552-86ae-0f0c7517bbda'
    )
ON CONFLICT DO NOTHING;

-- Assign Permission Groups to Super Admin Role
INSERT INTO "modules_roles" (
    "permission_group_id",
    "role_id"
)
VALUES
-- Backing Module to Super Admin Role Scope BEGIN
    (   
        '51dab038-4d4c-4d7f-8b94-2f3d642d13ed',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        '557777a3-022d-4013-985e-26c60a7e589a',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        'e13d1dbf-d235-45ff-8e54-b0ad2c8fbaa3',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        '59a0721f-08cb-4e12-a23b-bdbbab47abdc',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        '39cc8499-e573-404f-a301-f0806988b0e9',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    )
ON CONFLICT DO NOTHING;
-- Backing Module to Super Admin Role Scope END

