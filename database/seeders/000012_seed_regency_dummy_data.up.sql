-- Seed Dummy Data for Regency Module
-- 2 Provinces, 4 Cities (2 per province), 8 Districts (2 per city), 16 Subdistricts (2 per district)

-- Insert Provinces
INSERT INTO "province" ("id", "name", "created_at", "updated_at") 
VALUES 
    ('1a1b1c1d-1e1f-4a1b-1c1d-1e1f1a1b1c1d', 'Jawa Barat', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('2a2b2c2d-2e2f-4a2b-2c2d-2e2f2a2b2c2d', 'Jawa Tengah', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Insert Cities (2 per province)
INSERT INTO "city" ("id", "province_id", "name", "created_at", "updated_at")
VALUES
    -- Cities for Jawa Barat
    ('3a3b3c3d-3e3f-4a3b-3c3d-3e3f3a3b3c3d', '1a1b1c1d-1e1f-4a1b-1c1d-1e1f1a1b1c1d', 'Bandung', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('4a4b4c4d-4e4f-4a4b-4c4d-4e4f4a4b4c4d', '1a1b1c1d-1e1f-4a1b-1c1d-1e1f1a1b1c1d', 'Bekasi', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- Cities for Jawa Tengah
    ('5a5b5c5d-5e5f-4a5b-5c5d-5e5f5a5b5c5d', '2a2b2c2d-2e2f-4a2b-2c2d-2e2f2a2b2c2d', 'Semarang', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('6a6b6c6d-6e6f-4a6b-6c6d-6e6f6a6b6c6d', '2a2b2c2d-2e2f-4a2b-2c2d-2e2f2a2b2c2d', 'Surakarta', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Insert Districts (2 per city)
INSERT INTO "district" ("id", "city_id", "name", "created_at", "updated_at")
VALUES
    -- Districts for Bandung
    ('7a7b7c7d-7e7f-4a7b-7c7d-7e7f7a7b7c7d', '3a3b3c3d-3e3f-4a3b-3c3d-3e3f3a3b3c3d', 'Cicendo', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('8a8b8c8d-8e8f-4a8b-8c8d-8e8f8a8b8c8d', '3a3b3c3d-3e3f-4a3b-3c3d-3e3f3a3b3c3d', 'Cimahi', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- Districts for Bekasi
    ('9a9b9c9d-9e9f-4a9b-9c9d-9e9f9a9b9c9d', '4a4b4c4d-4e4f-4a4b-4c4d-4e4f4a4b4c4d', 'Bekasi Barat', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('0a0b0c0d-0e0f-4a0b-0c0d-0e0f0a0b0c0d', '4a4b4c4d-4e4f-4a4b-4c4d-4e4f4a4b4c4d', 'Bekasi Timur', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- Districts for Semarang
    ('a1b1c1d1-e1f1-4a1b-1c1d-1e1f1a1b1c1d', '5a5b5c5d-5e5f-4a5b-5c5d-5e5f5a5b5c5d', 'Semarang Barat', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('b2c2d2e2-f2a2-4b2c-2d2e-2f2a2b2c2d2e', '5a5b5c5d-5e5f-4a5b-5c5d-5e5f5a5b5c5d', 'Semarang Timur', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- Districts for Surakarta
    ('c3d3e3f3-a3b3-4c3d-3e3f-3a3b3c3d3e3f', '6a6b6c6d-6e6f-4a6b-6c6d-6e6f6a6b6c6d', 'Laweyan', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('d4e4f4a4-b4c4-4d4e-4f4a-4b4c4d4e4f4a', '6a6b6c6d-6e6f-4a6b-6c6d-6e6f6a6b6c6d', 'Serengan', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Insert Subdistricts (2 per district)
INSERT INTO "subdistrict" ("id", "district_id", "name", "created_at", "updated_at")
VALUES
    -- Subdistricts for Cicendo
    ('e5f5a5b5-c5d5-4e5f-5a5b-5c5d5e5f5a5b', '7a7b7c7d-7e7f-4a7b-7c7d-7e7f7a7b7c7d', 'Arjuna', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('f6a6b6c6-d6e6-4f6a-6b6c-6d6e6f6a6b6c', '7a7b7c7d-7e7f-4a7b-7c7d-7e7f7a7b7c7d', 'Pasteur', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- Subdistricts for Cimahi
    ('a7b7c7d7-e7f7-4a7b-7c7d-7e7f7a7b7c7d', '8a8b8c8d-8e8f-4a8b-8c8d-8e8f8a8b8c8d', 'Cibabat', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('b8c8d8e8-f8a8-4b8c-8d8e-8f8a8b8c8d8e', '8a8b8c8d-8e8f-4a8b-8c8d-8e8f8a8b8c8d', 'Cimahi Tengah', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- Subdistricts for Bekasi Barat
    ('c9d9e9f9-a9b9-4c9d-9e9f-9a9b9c9d9e9f', '9a9b9c9d-9e9f-4a9b-9c9d-9e9f9a9b9c9d', 'Bintara', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('d0e0f0a0-b0c0-4d0e-0f0a-0b0c0d0e0f0a', '9a9b9c9d-9e9f-4a9b-9c9d-9e9f9a9b9c9d', 'Jakasampurna', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- Subdistricts for Bekasi Timur
    ('e1f1a1b1-c1d1-4e1f-1a1b-1c1d1e1f1a1b', '0a0b0c0d-0e0f-4a0b-0c0d-0e0f0a0b0c0d', 'Aren Jaya', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('f2a2b2c2-d2e2-4f2a-2b2c-2d2e2f2a2b2c', '0a0b0c0d-0e0f-4a0b-0c0d-0e0f0a0b0c0d', 'Bekasi Jaya', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- Subdistricts for Semarang Barat
    ('a3b3c3d3-e3f3-4a3b-3c3d-3e3f3a3b3c3d', 'a1b1c1d1-e1f1-4a1b-1c1d-1e1f1a1b1c1d', 'Bongsari', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('b4c4d4e4-f4a4-4b4c-4d4e-4f4a4b4c4d4e', 'a1b1c1d1-e1f1-4a1b-1c1d-1e1f1a1b1c1d', 'Bojongsalaman', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- Subdistricts for Semarang Timur
    ('c5d5e5f5-a5b5-4c5d-5e5f-5a5b5c5d5e5f', 'b2c2d2e2-f2a2-4b2c-2d2e-2f2a2b2c2d2e', 'Bugangan', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('d6e6f6a6-b6c6-4d6e-6f6a-6b6c6d6e6f6a', 'b2c2d2e2-f2a2-4b2c-2d2e-2f2a2b2c2d2e', 'Karangtempel', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- Subdistricts for Laweyan
    ('e7f7a7b7-c7d7-4e7f-7a7b-7c7d7e7f7a7b', 'c3d3e3f3-a3b3-4c3d-3e3f-3a3b3c3d3e3f', 'Bumi', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('f8a8b8c8-d8e8-4f8a-8b8c-8d8e8f8a8b8c', 'c3d3e3f3-a3b3-4c3d-3e3f-3a3b3c3d3e3f', 'Kestalan', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- Subdistricts for Serengan
    ('a9b9c9d9-e9f9-4a9b-9c9d-9e9f9a9b9c9d', 'd4e4f4a4-b4c4-4d4e-4f4a-4b4c4d4e4f4a', 'Gandekan', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('b0c0d0e0-f0a0-4b0c-0d0e-0f0a0b0c0d0e', 'd4e4f4a4-b4c4-4d4e-4f4a-4b4c4d4e4f4a', 'Kepatihan Kulon', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;
