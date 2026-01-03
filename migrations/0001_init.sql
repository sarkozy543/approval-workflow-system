CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS approval_requests (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_env     TEXT NOT NULL,
    target_env     TEXT NOT NULL,
    change_payload JSONB NOT NULL,
    status         TEXT NOT NULL DEFAULT 'PENDING',
    requested_by   TEXT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    repo_url       TEXT,
    source_commit  TEXT,
    target_commit  TEXT,
    change_type    TEXT,
    source_branch  TEXT
);

CREATE TABLE IF NOT EXISTS approval_logs (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id uuid NOT NULL REFERENCES approval_requests(id) ON DELETE CASCADE,
    action     TEXT NOT NULL, -- created / approved / rejected
    action_by  TEXT NOT NULL,
    note       TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
