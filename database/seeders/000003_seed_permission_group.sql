INSERT INTO "permission_groups" ("id", "created_at", "updated_at", "name", "deletable", "description", "module") 
VALUES 
-- User
    ('b72e75d3-73fe-43fc-8b24-06b74f5c707d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Add User ', false, 'Have Full Access for Add User Sub-Module', 'Users'),
    ('fa2bcd6c-e4c7-4bda-b80e-cf9d3fdac557', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update User', false, 'Have Full Access for Update User Sub-Module', 'Users'),
    ('a51a23c9-dab9-4b61-b38e-42b52c68bb57', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Block User', false, 'Have Full Access for Block User Sub-Module', 'Users'),
    ('5f1c4808-6c8d-445b-bab1-2d1704531ff5', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'View User', false, 'Have Full Access for View User Sub-Module', 'Users'),
    ('f9a3bff0-22f9-45f5-9385-623206b5b4da', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Add Role', false, 'Have Full Access for Add Role Sub-Module', 'Roles'),
    ('3a02c1a9-d4c4-42a4-9f52-0c3b6794b580', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Update Role', false, 'Have Full Access for Update Role Sub-Module', 'Roles'),
    ('d2193e52-4e89-4b50-b7ab-745d6ab36a22', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Delete Role', false, 'Have Full Access for Delete Role Sub-Module', 'Roles'),
    ('15a52416-c65a-4155-80a8-160602fd22fe', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'View Role', false, 'Have Full Access for View Role Sub-Module', 'Roles');
