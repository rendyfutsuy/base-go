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
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        'Super Admin',
        'false'
    ),
    (
       '0199dc4e-6455-7b8d-b48b-409dbece678b',
        '2023-09-25 15:33:01.881',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        'User',
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