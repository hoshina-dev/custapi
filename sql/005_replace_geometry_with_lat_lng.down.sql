CREATE EXTENSION IF NOT EXISTS postgis;

DROP INDEX IF EXISTS organizations_lat_lng_idx;

ALTER TABLE organizations
ADD COLUMN geom GEOMETRY(POINT, 4326);

UPDATE organizations
SET geom = ST_Point(longitude, latitude, 4326)
WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

ALTER TABLE organizations
ALTER COLUMN geom SET NOT NULL;

CREATE INDEX IF NOT EXISTS organizations_geom_idx ON organizations USING GIST (geom);

ALTER TABLE organizations
DROP COLUMN latitude,
DROP COLUMN longitude;
