ALTER TABLE users DROP CONSTRAINT fk_organization;

ALTER TABLE users 
ADD CONSTRAINT fk_organization 
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

ALTER TABLE users
DROP COLUMN IF EXISTS password,
DROP COLUMN IF EXISTS is_admin,
DROP COLUMN IF EXISTS phone_number,
DROP COLUMN IF EXISTS social_media,
DROP COLUMN IF EXISTS description, 
DROP COLUMN IF EXISTS avatar_url,
DROP COLUMN IF EXISTS research_categories;