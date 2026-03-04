-- what do we do with soft delete? soft delete flag -> cascade to delete here?
INSERT INTO user_organizations (user_id, organization_id, is_admin) 
SELECT u.id, u.organization_id, COALESCE(u.is_admin, FALSE) 
FROM users u
WHERE u.organization_id IS NOT NULL 
AND deleted_at IS NULL 
ON CONFLICT (user_id, organization_id) DO NOTHING;
