ALTER TABLE users
    DROP CONSTRAINT IF EXISTS fk_organization,
    DROP COLUMN IF EXISTS organization_id,
    -- i thhink admin at user org level make more sense
    DROP COLUMN IF EXISTS is_admin;
