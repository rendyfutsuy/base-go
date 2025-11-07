-- Remove unique constraint from email column in users table
-- PostgreSQL automatically creates a constraint named 'users_email_key' when UNIQUE is specified in column definition
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;

