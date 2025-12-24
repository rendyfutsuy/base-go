-- Seed Parameters for Supplier Module
INSERT INTO "parameters" (
    "id",
    "code",
    "name",
    "value",
    "type",
    "description",
    "created_at",
    "updated_at",
    "deleted_at"
)
VALUES
-- delivery_option
    (uuid_generate_v7(), 'SUPPLIER_DELIVERY_OPTION_SENT_BY_SUPPLIER', 'Dikirim Supplier', 'sent_by_supplier', 'delivery_option', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'SUPPLIER_DELIVERY_OPTION_PICK_UP_BY_CUSTOMER', 'Diambil Sendiri', 'pick_up_by_customer', 'delivery_option', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'SUPPLIER_DELIVERY_OPTION_SENT_BY_COURIER', 'Ekspedisi', 'sent_by_courier', 'delivery_option', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),

-- expedition_paid_by
    (uuid_generate_v7(), 'SUPPLIER_BOUGHT_BY_PAID_BY_SUPPLIER', 'Bayar Supplier', 'paid_by_supplier', 'expedition_paid_by', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'SUPPLIER_BOUGHT_BY_PAID_BY_CUSTOMER', 'Bayar Sendiri', 'paid_by_customer', 'expedition_paid_by', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),

-- expedition_calculation
    (uuid_generate_v7(), 'SUPPLIER_CALCULATION_DELIVERY_TYPE_FULL', 'Kalkulasi Faktur (100%)', 'full', 'expedition_calculation', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'SUPPLIER_CALCULATION_DELIVERY_TYPE_HALF', 'Kalkulasi Faktur (50%)', 'half', 'expedition_calculation', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    
-- identity_type
    (uuid_generate_v7(), 'IDENTITY_TYPE_KTP', 'KTP', 'ktp', 'identity_type', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'IDENTITY_TYPE_NPWP', 'NPWP', 'npwp', 'identity_type', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'IDENTITY_TYPE_OTHER', 'Identitas Lainnya', 'others', 'identity_type', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL)
ON CONFLICT (code) DO NOTHING;

