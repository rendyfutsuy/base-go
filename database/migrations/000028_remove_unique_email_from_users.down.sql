-- Restore unique constraint on email column in users table
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);

