INSERT INTO
    "roles" (
        "id",
        "created_at",
        "created_by",
        "updated_at",
        "updated_by",
        "deleted_at",
        "deleted_by",
        "name",
        "deletable"
    )
VALUES
    (
        'a43a5e5f-a172-42d1-a70e-8834bf653eb0',
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Super Admin',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Technical Support',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Under Writer',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Business Support Officer',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Business Support Manager',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Under Writer Treaty',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Group Head Treaty',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'User Preview',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Claim Analyst',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Senior Claim Analyst',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Claim Business Analyst',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Group Head Claim',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Claim and Admin Group Manager',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'Technical Director',
        'false'
    ),
    (
        uuid_generate_v7(),
        '2023-09-25 15:33:01.881',
        'system',
        NULL,
        NULL,
        NULL,
        NULL,
        'ser Preview Claim Analyst',
        'false'
    )
     ON CONFLICT ("id") DO
UPDATE
SET
    "created_at" = EXCLUDED."created_at",
    "created_by" = EXCLUDED."created_by",
    "updated_at" = EXCLUDED."updated_at",
    "updated_by" = EXCLUDED."updated_by",
    "deleted_at" = EXCLUDED."deleted_at",
    "deleted_by" = EXCLUDED."deleted_by",
    "name" = EXCLUDED."name",
    "deletable" = EXCLUDED."deletable";