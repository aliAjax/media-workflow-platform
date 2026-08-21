CREATE TABLE IF NOT EXISTS media_assets (id text PRIMARY KEY, tenant_id text NOT NULL, name text NOT NULL, state text NOT NULL, created_at timestamptz NOT NULL, updated_at timestamptz NOT NULL);
CREATE TABLE IF NOT EXISTS media_jobs (id text PRIMARY KEY, tenant_id text NOT NULL, asset_id text NOT NULL, pipeline_id text NOT NULL, state text NOT NULL, progress integer NOT NULL, attempt integer NOT NULL, created_at timestamptz NOT NULL, updated_at timestamptz NOT NULL);
CREATE INDEX IF NOT EXISTS media_jobs_tenant_state ON media_jobs(tenant_id,state);
