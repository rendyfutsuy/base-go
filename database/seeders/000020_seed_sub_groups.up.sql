-- Seed Data for Sub Groups Table
-- UUID static yang valid dan unique
-- Reference ke group yang sudah di-seed di 000023_seed_group.up.sql

INSERT INTO "sub_groups" ("id", "groups_id", "name", "created_at", "updated_at", "deleted_at", "deleted_by") 
VALUES 
    -- Sub Groups untuk Elektronik (a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d)
    ('1a2b3c4d-5e6f-4a7b-8c9d-0e1f2a3b4c5d', 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'SMARTPHONE', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('2b3c4d5e-6f7a-4b8c-9d0e-1f2a3b4c5d6e', 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'LAPTOP', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('3c4d5e6f-7a8b-4c9d-0e1f-2a3b4c5d6e7f', 'a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'TV & AUDIO', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    
    -- Sub Groups untuk Fashion (b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e)
    ('4d5e6f7a-8b9c-4d0e-1f2a-3b4c5d6e7f8a', 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'PAKAIAN PRIA', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('5e6f7a8b-9c0d-4e1f-2a3b-4c5d6e7f8a9b', 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'PAKAIAN WANITA', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('6f7a8b9c-0d1e-4f2a-3b4c-5d6e7f8a9b0c', 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'SEPATU & SANDAL', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    
    -- Sub Groups untuk Makanan & Minuman (c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f)
    ('7a8b9c0d-1e2f-4a3b-4c5d-6e7f8a9b0c1d', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'MAKANAN RINGAN', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('8b9c0d1e-2f3a-4b4c-5d6e-7f8a9b0c1d2e', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'MINUMAN', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('9c0d1e2f-3a4b-4c5d-6e7f-8a9b0c1d2e3f', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'BAHAN MAKANAN', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    
    -- Sub Groups untuk Kesehatan & Kecantikan (d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a)
    ('0d1e2f3a-4b5c-4d6e-7f8a-9b0c1d2e3f4a', 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'SKINCARE', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('1e2f3a4b-5c6d-4e7f-8a9b-0c1d2e3f4a5b', 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'OBAT-OBATAN', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('2f3a4b5c-6d7e-4f8a-9b0c-1d2e3f4a5b6c', 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'VITAMIN & SUPLEMEN', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    
    -- Sub Groups untuk Olahraga & Outdoor (e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b)
    ('3a4b5c6d-7e8f-4a9b-0c1d-2e3f4a5b6c7d', 'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'ALAT OLAHRAGA', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('4b5c6d7e-8f9a-4b0c-1d2e-3f4a5b6c7d8e', 'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'PAKAIAN OLAHRAGA', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    
    -- Sub Groups untuk Buku & Alat Tulis (f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c)
    ('5c6d7e8f-9a0b-4c1d-2e3f-4a5b6c7d8e9f', 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'BUKU FIKSI', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('6d7e8f9a-0b1c-4d2e-3f4a-5b6c7d8e9f0a', 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'BUKU NON-FIKSI', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('7e8f9a0b-1c2d-4e3f-4a5b-6c7d8e9f0a1b', 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'ALAT TULIS', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    
    -- Sub Groups untuk Mainan & Hobi (a7b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d)
    ('8f9a0b1c-2d3e-4f4a-5b6c-7d8e9f0a1b2c', 'a7b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', 'MAINAN ANAK', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('9a0b1c2d-3e4f-4a5b-6c7d-8e9f0a1b2c3d', 'a7b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', 'HOBI & KOLEKSI', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    
    -- Sub Groups untuk Otomotif (b8c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e)
    ('0b1c2d3e-4f5a-4b6c-7d8e-9f0a1b2c3d4e', 'b8c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', 'SPAREPART MOTOR', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('1c2d3e4f-5a6b-4c7d-8e9f-0a1b2c3d4e5f', 'b8c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', 'SPAREPART MOBIL', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('2d3e4f5a-6b7c-4d8e-9f0a-1b2c3d4e5f6a', 'b8c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', 'AKSESORIS OTOMOTIF', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    
    -- Sub Groups untuk Rumah Tangga (c9d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f)
    ('3e4f5a6b-7c8d-4e9f-0a1b-2c3d4e5f6a7b', 'c9d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f', 'PERABOTAN', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('4f5a6b7c-8d9e-4f0a-1b2c-3d4e5f6a7b8c', 'c9d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f', 'DAPUR & MAKANAN', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('5a6b7c8d-9e0f-4a1b-2c3d-4e5f6a7b8c9d', 'c9d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f', 'KAMAR MANDI', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    
    -- Sub Groups untuk Pertanian & Perkebunan (d0e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a)
    ('6b7c8d9e-0f1a-4b2c-3d4e-5f6a7b8c9d0e', 'd0e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a', 'BENIH & BIBIT', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('7c8d9e0f-1a2b-4c3d-4e5f-6a7b8c9d0e1f', 'd0e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a', 'PUPUK & PESTISIDA', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL),
    ('8d9e0f1a-2b3c-4d4e-5f6a-7b8c9d0e1f2a', 'd0e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a', 'ALAT PERTANIAN', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, NULL)
ON CONFLICT (id) DO NOTHING;

