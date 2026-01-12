ALTER TABLE organizations
DROP COLUMN IF EXISTS geom,
DROP COLUMN IF EXISTS address,
DROP COLUMN IF EXISTS description,
DROP COLUMN IF EXISTS image_urls;

DROP INDEX IF EXISTS organizations_geom_idx;