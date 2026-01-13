ALTER TABLE users DROP CONSTRAINT fk_organization;

ALTER TABLE users 
ADD CONSTRAINT fk_organization 
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE RESTRICT;

ALTER TABLE users
ADD COLUMN password TEXT NOT NULL,
ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN phone_number VARCHAR(25),
ADD COLUMN social_media TEXT,
ADD COLUMN description TEXT, 
ADD COLUMN avatar_url TEXT,
ADD COLUMN research_categories TEXT[] NOT NULL DEFAULT '{}';