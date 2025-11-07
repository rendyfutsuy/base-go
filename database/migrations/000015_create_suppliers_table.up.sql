-- Create sequence for supplier_code starting from 1
CREATE SEQUENCE IF NOT EXISTS supplier_code_seq START WITH 1 INCREMENT BY 1;

-- Function to generate formatted supplier_code: "0" + seq if 1 digit, else just seq
CREATE OR REPLACE FUNCTION generate_supplier_code()
RETURNS VARCHAR AS $$
DECLARE
  seq_val BIGINT;
  formatted_code VARCHAR;
BEGIN
  seq_val := nextval('supplier_code_seq');
  IF seq_val < 10 THEN
    formatted_code := '0' || seq_val::VARCHAR;
  ELSE
    formatted_code := seq_val::VARCHAR;
  END IF;
  RETURN formatted_code;
END;
$$ LANGUAGE plpgsql;

-- Create table suppliers
CREATE TABLE IF NOT EXISTS suppliers (
  id UUID DEFAULT uuid_generate_v7() PRIMARY KEY NOT NULL,
  
  subdistrict_id UUID REFERENCES subdistrict(id),
  expedition_arrives_id UUID REFERENCES expeditions(id),
  
  created_at TIMESTAMP,
  created_by VARCHAR(255),
  updated_at TIMESTAMP,
  updated_by VARCHAR(255),
  deleted_at TIMESTAMP,
  deleted_by VARCHAR(255),
  
  supplier_code VARCHAR(255) NOT NULL UNIQUE DEFAULT generate_supplier_code(),
  
  identity_type VARCHAR(255) NOT NULL,
  identity_name VARCHAR(255),
  identity_number VARCHAR(255) NOT NULL UNIQUE,
  identity_document TEXT,
  
  supplier_name VARCHAR(255) NOT NULL,
  alias VARCHAR(255),
  rt VARCHAR(255),
  rw VARCHAR(255),
  postal_code VARCHAR(255),
  address TEXT NOT NULL,
  
  telp_number VARCHAR(255),
  phone_number VARCHAR(255),
  email VARCHAR(255),
  notes TEXT,
  
  relation_date DATE DEFAULT CURRENT_DATE,
  delivery_option UUID NOT NULL REFERENCES parameter(id),
  expedition_paid_by UUID REFERENCES parameter(id),
  expedition_calculation UUID REFERENCES parameter(id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS suppliers_id_index ON suppliers (id);
CREATE INDEX IF NOT EXISTS suppliers_supplier_code_index ON suppliers (supplier_code);
CREATE INDEX IF NOT EXISTS suppliers_supplier_name_index ON suppliers (supplier_name);
CREATE INDEX IF NOT EXISTS suppliers_identity_number_index ON suppliers (identity_number);
CREATE INDEX IF NOT EXISTS suppliers_subdistrict_id_index ON suppliers (subdistrict_id);
CREATE INDEX IF NOT EXISTS suppliers_expedition_arrives_id_index ON suppliers (expedition_arrives_id);
CREATE INDEX IF NOT EXISTS suppliers_delivery_option_index ON suppliers (delivery_option);
CREATE INDEX IF NOT EXISTS suppliers_identity_type_index ON suppliers (identity_type);
CREATE INDEX IF NOT EXISTS suppliers_email_index ON suppliers (email);
CREATE INDEX IF NOT EXISTS suppliers_created_at_index ON suppliers (created_at);
CREATE INDEX IF NOT EXISTS suppliers_updated_at_index ON suppliers (updated_at);
CREATE INDEX IF NOT EXISTS suppliers_deleted_at_index ON suppliers (deleted_at);

