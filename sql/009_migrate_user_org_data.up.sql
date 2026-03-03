-- what do we do with soft delete? soft delete flag -> cascade to delete here?
INSERT INTO user_organization (user_id, organization_id, is_admin) 
SELECT id, organization_id, COALESCE(is_admin, FALSE) 
FROM users 
WHERE organization_id IS NOT NULL 
AND deleted_at IS NULL 
ON CONFLICT (user_id, organization_id) DO NOTHING;
