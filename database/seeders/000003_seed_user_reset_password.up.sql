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
    -- Super User - Buat Coba Coba Reset Password
    -- password: password123
    -- role: Super Admin
    -- password expired at: now + 90 days
    -- status active
    uuid_generate_v7(),
    NOW(),
    NOW(),
    NOW() + INTERVAL '90 days',
    'superuser@mailinator.com',
    'admin',
    '$2y$10$GViZu3GONfwoswHMagB0sOh.ZlKeK9WrSyhwbvmiheeGGihz2vBSm',
    'a43a5e5f-a172-42d1-a70e-8834bf653eb0',
    'Super User - Buat Coba Coba Reset Password',
    TRUE,
    'male',
    0,
    '3276052208000023'
),
(
    -- Super User - Password Expired
    -- password: password123
    -- role: Super Admin
    -- password expired at: now - 10 days
    -- status active
    uuid_generate_v7(),
    NOW(),
    NOW(),
    NOW() - INTERVAL '10 days',
    'superuser.password.expired@mailinator.com',
    'admin',
    '$2y$10$GViZu3GONfwoswHMagB0sOh.ZlKeK9WrSyhwbvmiheeGGihz2vBSm',
    'a43a5e5f-a172-42d1-a70e-8834bf653eb0',
    'Super User - Password Expired',
    TRUE,
    'male',
    0,
    '3276052223000032'
);
