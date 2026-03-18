-- Rename 'admin' role to 'manager' in user_organizations (organization-level role)
UPDATE user_organizations
SET role = 'manager'
WHERE role = 'admin';

-- Add site-level role column to users table
ALTER TABLE users
ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'user';
