-- Migration: 001_create_organizations_table
-- Description: Rollback organizations table creation

DROP TRIGGER IF EXISTS update_organizations_updated_at ON organizations;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_organizations_deleted_at;
DROP TABLE IF EXISTS organizations;
