-- Seed Permission Groups for Module "Course"
INSERT INTO "permission_groups" ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
    ('c1a9a6e8-7c3f-4c9a-9b6e-1f2d3e4a5b6c', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Create', false, 'Have Full Access for Create Course Sub-Module', 'Post'),
    ('b2b8a7d9-8d4e-4e0b-8c7f-2e3f4a5b6c7d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update', false, 'Have Full Access for Update Course Sub-Module', 'Post'),
    ('a3c7b8e9-9e5f-4f1c-9d8e-3f4a5b6c7d8e', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete', false, 'Have Full Access for Delete Course Sub-Module', 'Post')
ON CONFLICT (id) DO NOTHING;

-- Seed Permissions for Module "Course"
INSERT INTO "permissions" ("id", "created_at", "updated_at", "name", "deletable")
VALUES
    ('d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'post.create', false),
    ('e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'post.update', false),
    ('f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'post.delete', false)
ON CONFLICT (id) DO NOTHING;

-- Map Permission Groups <-> Permissions for Module "Course"
INSERT INTO "permissions_modules" ("permission_group_id", "permission_id")
VALUES
    ('c1a9a6e8-7c3f-4c9a-9b6e-1f2d3e4a5b6c', 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a'),
    ('b2b8a7d9-8d4e-4e0b-8c7f-2e3f4a5b6c7d', 'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b'),
    ('a3c7b8e9-9e5f-4f1c-9d8e-3f4a5b6c7d8e', 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c')
ON CONFLICT DO NOTHING;

-- Assign Permission Groups to Super Admin Role
INSERT INTO "modules_roles" (
    "permission_group_id",
    "role_id"
)
VALUES
    (   
        'c1a9a6e8-7c3f-4c9a-9b6e-1f2d3e4a5b6c',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        'b2b8a7d9-8d4e-4e0b-8c7f-2e3f4a5b6c7d',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    ),
    (   
        'a3c7b8e9-9e5f-4f1c-9d8e-3f4a5b6c7d8e',
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0'
    )
ON CONFLICT DO NOTHING;
