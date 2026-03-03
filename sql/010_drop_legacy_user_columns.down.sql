ALTER TABLE users
    ADD COLUMN IF NOT EXISTS organization_id UUID;
    ADD COLUMN IF NOT EXISTS is_admin BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE users u
SET organization_id = uo.organization_id
FROM (
    SELECT DISTINCT ON (user_id)
        user_id,
        organization_id
    FROM  user_organizations
    ORDER BY user_id, created_at ASC
) uo
WHERE u.id = uo.user_id;

UPDATE users u
SET is_admin = TRUE
FROM user_organizations uo
WHERE uo.user_id = u.id
AND uo.is_admin = TRUE;

ALTER TABLE users
    ADD CONSTRAINT fk_organization
    FOREIGN KEY (organization_id)
    REFERENCES organizations(id)
    ON DELETE RESTRICT;
