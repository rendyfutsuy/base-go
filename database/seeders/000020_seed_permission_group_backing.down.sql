-- Remove Permission Groups assignments from Super Admin Role
DELETE FROM "modules_roles" 
WHERE "permission_group_id" IN (
    '51dab038-4d4c-4d7f-8b94-2f3d642d13ed',
    '557777a3-022d-4013-985e-26c60a7e589a',
    'e13d1dbf-d235-45ff-8e54-b0ad2c8fbaa3',
    '59a0721f-08cb-4e12-a23b-bdbbab47abdc',
    '39cc8499-e573-404f-a301-f0806988b0e9'
);

-- Remove Permissions Modules mappings
DELETE FROM "permissions_modules" 
WHERE "permission_group_id" IN (
    '51dab038-4d4c-4d7f-8b94-2f3d642d13ed',
    '557777a3-022d-4013-985e-26c60a7e589a',
    'e13d1dbf-d235-45ff-8e54-b0ad2c8fbaa3',
    '59a0721f-08cb-4e12-a23b-bdbbab47abdc',
    '39cc8499-e573-404f-a301-f0806988b0e9'
);

-- Remove Permissions
DELETE FROM "permissions" 
WHERE "id" IN (
    '2ba6c0e3-c935-4856-bcac-4fe56b8d975d',
    '1a27a520-b66d-4b84-b67c-18b009cb8754',
    '9c5b5801-7602-4631-bbee-86a878d414c9',
    'c69a9c7e-aac9-4463-8a9a-04f1ecc55a02',
    'ed2c4b27-f1a7-4552-86ae-0f0c7517bbda'
);

-- Remove Permission Groups
DELETE FROM "permission_groups" 
WHERE "id" IN (
    '51dab038-4d4c-4d7f-8b94-2f3d642d13ed',
    '557777a3-022d-4013-985e-26c60a7e589a',
    'e13d1dbf-d235-45ff-8e54-b0ad2c8fbaa3',
    '59a0721f-08cb-4e12-a23b-bdbbab47abdc',
    '39cc8499-e573-404f-a301-f0806988b0e9'
);

