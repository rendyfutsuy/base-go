INSERT INTO users (
    id,
    created_at,
    updated_at,
    password_expired_at,
    email,
    username,
    password,
    role_id,
    full_name,
    is_active,
    gender,
    counter,
    nik
) VALUES
(
    -- Super Admin User
    -- password: password123
    -- role: Super Admin
    -- password expired at: now + 90 days
    -- status active
    uuid_generate_v7(),
    NOW(),
    NOW(),
    NOW() + INTERVAL '90 days',
    'admin@mailinator.com',
    'admin',
    '$2y$10$GViZu3GONfwoswHMagB0sOh.ZlKeK9WrSyhwbvmiheeGGihz2vBSm',
    'a43a5e5f-a172-42d1-a70e-8834bf653eb0',
    'Admin - System',
    TRUE,
    'male',
    0,
    '3175091206980001'
),
(
    -- Super Admin User
    -- password: password123
    -- role: Super Admin
    -- password expired at: now - 2 days
    -- status active
    uuid_generate_v7(),
    NOW(),
    NOW(),
    NOW() - INTERVAL '2 days',
    'admin.password.expired@mailinator.com',
    'admin.password.expired',
    '$2y$10$GViZu3GONfwoswHMagB0sOh.ZlKeK9WrSyhwbvmiheeGGihz2vBSm',
    'a43a5e5f-a172-42d1-a70e-8834bf653eb0',
    'Admin - Password Expired',
    TRUE,
    'male',
    0,
    '3276012501800002'
),
(
    -- Super ADmin User
    -- password: password123
    -- role: Super Admin
    -- password expired at: now + 90 days
    -- status inactive
    uuid_generate_v7(),
    NOW(),
    NOW(),
    NOW() + INTERVAL '90 days',
    'admin.inactive@mailinator.com',
    'admin.inactive',
    '$2y$10$GViZu3GONfwoswHMagB0sOh.ZlKeK9WrSyhwbvmiheeGGihz2vBSm',
    'a43a5e5f-a172-42d1-a70e-8834bf653eb0',
    'Admin - Inactive',
    FALSE,
    'male',
    0,
    '3374121301900003'
);
