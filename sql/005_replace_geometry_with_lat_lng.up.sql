ALTER TABLE organizations
ADD COLUMN latitude NUMERIC(10, 8),
ADD COLUMN longitude NUMERIC(11, 8);

UPDATE organizations
SET latitude = ST_Y(geom),
    longitude = ST_X(geom)
WHERE geom IS NOT NULL;

ALTER TABLE organizations
ALTER COLUMN latitude SET NOT NULL,
ALTER COLUMN longitude SET NOT NULL;

DROP INDEX IF EXISTS organizations_geom_idx;
ALTER TABLE organizations DROP COLUMN geom;

CREATE INDEX IF NOT EXISTS organizations_lat_lng_idx ON organizations (latitude, longitude);

DROP EXTENSION IF EXISTS postgis;