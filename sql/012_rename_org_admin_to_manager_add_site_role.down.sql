-- Drop site-level role column from users table
ALTER TABLE users
DROP COLUMN role;

-- Revert 'manager' back to 'admin' in user_organizations
UPDATE user_organizations
SET role = 'admin'
WHERE role = 'manager';
