-- Migration: 007_add_trigram_indexes
CREATE INDEX IF NOT EXISTS idx_users_name_trgm ON users USING GIN (name gin_trgm_ops);
-- Create trigram index on users.name

CREATE INDEX IF NOT EXISTS idx_organizations_name_trgm ON organizations USING GIN (name gin_trgm_ops);
-- Create trigram index on organizations.name

-- Description: Add trigram indexes on users and organizations name columns for fuzzy search

