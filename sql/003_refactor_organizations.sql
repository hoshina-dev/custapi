ALTER TABLE organizations
ADD COLUMN geom GEOMETRY(POINT, 4326) NOT NULL,
ADD COLUMN address TEXT,
ADD COLUMN description TEXT,
ADD COLUMN image_urls TEXT[];

CREATE INDEX IF NOT EXIST organizations_geom_idx ON organizations USING GIST (geom);