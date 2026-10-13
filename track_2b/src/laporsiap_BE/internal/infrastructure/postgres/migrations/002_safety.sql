ALTER TABLE cases ADD COLUMN IF NOT EXISTS expires_at timestamptz;
ALTER TABLE cases ADD COLUMN IF NOT EXISTS access_expires_at timestamptz;
UPDATE cases SET expires_at = created_at + interval '30 days' WHERE expires_at IS NULL;
UPDATE cases SET access_expires_at = LEAST(expires_at, created_at + interval '7 days') WHERE access_expires_at IS NULL;
ALTER TABLE cases ALTER COLUMN expires_at SET NOT NULL;
ALTER TABLE cases ALTER COLUMN access_expires_at SET NOT NULL;
CREATE INDEX IF NOT EXISTS cases_expiry_idx ON cases(expires_at);
CREATE TABLE IF NOT EXISTS ai_attempts (
 id text PRIMARY KEY,
 case_id text REFERENCES cases(id) ON DELETE SET NULL,
 created_at timestamptz NOT NULL DEFAULT now(),
 lease_until timestamptz NOT NULL DEFAULT now() + interval '45 seconds',
 finished_at timestamptz
);
CREATE INDEX IF NOT EXISTS ai_attempts_created_idx ON ai_attempts(created_at);
CREATE INDEX IF NOT EXISTS ai_attempts_case_idx ON ai_attempts(case_id, created_at);
INSERT INTO schema_migrations(version) VALUES (2) ON CONFLICT DO NOTHING;
