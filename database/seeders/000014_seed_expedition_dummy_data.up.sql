-- Seed Dummy Data for Expedition Module (4 records)
INSERT INTO "expeditions" ("id", "expedition_code", "expedition_name", "address", "notes", "created_at", "created_by", "updated_at", "updated_by") 
VALUES 
    ('1a2b3c4d-5e6f-4a7b-8c9d-0e1f2a3b4c5d', '01', 'JNE', 'Jl. Ahmad Yani No. 123', 'Ekspedisi cepat dan terpercaya', CURRENT_TIMESTAMP, 'system', CURRENT_TIMESTAMP, 'system'),
    ('2b3c4d5e-6f7a-4b8c-9d0e-1f2a3b4c5d6e', '02', 'TIKI', 'Jl. Sudirman No. 456','Pengiriman seluruh Indonesia', CURRENT_TIMESTAMP, 'system', CURRENT_TIMESTAMP, 'system'),
    ('3c4d5e6f-7a8b-4c9d-0e1f-2a3b4c5d6e7f', '03', 'J&T Express', 'Jl. Gatot Subroto No. 789', 'Ekspedisi murah dan cepat', CURRENT_TIMESTAMP, 'system', CURRENT_TIMESTAMP, 'system'),
    ('4d5e6f7a-8b9c-4d0e-1f2a-3b4c5d6e7f8a', '04', 'SiCepat', 'Jl. Thamrin No. 321', 'Pengiriman door to door', CURRENT_TIMESTAMP, 'system', CURRENT_TIMESTAMP, 'system')
ON CONFLICT (id) DO NOTHING;

