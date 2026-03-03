UPDATE users u
SET organization_id = uo.organization_id
FROM (
    SELECT DISTINCT ON (user_id)
        user_id,
        organization_id
    FROM  user_organization
    ORDER BY user_id, created_at ASC
) uo
WHERE u.id = uo.user_id
AND u.organization_id IS NULL;

UPDATE users u
SET is_admin = EXISTS (
    SELECT 1
    FROM   user_organization uo
    WHERE  uo.user_id = u.id
    AND  uo.is_admin = TRUE
);

DELETE FROM user_organization;
