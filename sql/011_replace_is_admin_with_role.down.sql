ALTER TABLE user_organizations
ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE user_organizations
SET is_admin = (role = 'admin');

ALTER TABLE user_organizations
DROP COLUMN role;
