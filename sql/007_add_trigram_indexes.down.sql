-- Migration: 007_add_trigram_indexes
-- Description: Rollback trigram indexes on users and organizations name columns

-- Drop trigram index on users.email
DROP INDEX IF EXISTS idx_users_email_trgm;

-- Drop trigram index on users.name
DROP INDEX IF EXISTS idx_users_name_trgm;

-- Drop trigram index on organizations.name
DROP INDEX IF EXISTS idx_organizations_name_trgm;

