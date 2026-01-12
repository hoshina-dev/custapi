-- Migration: 002_create_users_table
-- Description: Rollback users table creation

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_organization_id;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
