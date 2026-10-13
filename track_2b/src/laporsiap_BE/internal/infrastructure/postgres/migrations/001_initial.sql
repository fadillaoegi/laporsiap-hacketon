CREATE TABLE IF NOT EXISTS schema_migrations (
  version integer PRIMARY KEY,
  applied_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS cases (
  id text PRIMARY KEY,
  token_hash bytea NOT NULL,
  revision bigint NOT NULL DEFAULT 1,
  document jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS case_history (
  case_id text NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
  revision bigint NOT NULL,
  action text NOT NULL,
  document jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (case_id, revision)
);
CREATE TABLE IF NOT EXISTS attachments (
  id text PRIMARY KEY,
  case_id text NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
  filename text NOT NULL,
  content_type text NOT NULL,
  sha256 text NOT NULL,
  data bytea NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS attachments_case_idx ON attachments(case_id);
INSERT INTO schema_migrations(version) VALUES (1) ON CONFLICT DO NOTHING;
