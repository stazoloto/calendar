CREATE TABLE IF NOT EXISTS events(
    id          BIGSERIAL       PRIMARY KEY,
    title       VARCHAR(255)    NOT NULL,
    description TEXT            NOT NULL DEFAULT '',
    start_at    TIMESTAMPTZ     NOT NULL,
    end_at      TIMESTAMPTZ     NOT NULL,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_events_start_at  ON events (start_at);
CREATE INDEX IF NOT EXISTS idx_events_end_at    ON events (end_at);