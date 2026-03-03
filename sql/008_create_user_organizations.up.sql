CREATE TABLE IF NOT EXISTS user_organization (
    user_id         UUID      NOT NULL,
    organization_id UUID      NOT NULL,
    is_admin        BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, organization_id),

    CONSTRAINT fk_uo_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    CONSTRAINT fk_uo_organization
    FOREIGN KEY (organization_id)
    REFERENCES organizations(id)
    ON DELETE CASCADE
);

-- Index organization id
CREATE INDEX IF NOT EXISTS idx_uo_org_id
ON user_organization(organization_id);

-- Partial Index for finding admin in an organization, could be useful and size is small cuz partial only admin
CREATE INDEX IF NOT EXISTS idx_uo_org_admin
ON user_organization(organization_id)
WHERE is_admin = TRUE;