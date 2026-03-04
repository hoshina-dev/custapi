ALTER TABLE user_organizations
ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'user';

UPDATE user_organizations
SET role = CASE WHEN is_admin THEN 'admin' ELSE 'user' END;

ALTER TABLE user_organizations
DROP COLUMN is_admin;
