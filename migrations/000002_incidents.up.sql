CREATE TABLE IF NOT EXISTS machine_incidents (
    id          VARCHAR(64) PRIMARY KEY,
    machine_id  VARCHAR(64) NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
    event_id    VARCHAR(64) NOT NULL,
    temp_c      DOUBLE PRECISION NOT NULL,
    vibro       DOUBLE PRECISION NOT NULL,
    reason      VARCHAR(255) NOT NULL,
    is_resolved BOOLEAN NOT NULL DEFAULT FALSE,
    
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_incidents_machine_id ON machine_incidents(machine_id, is_resolved);
