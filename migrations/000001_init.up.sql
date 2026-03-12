CREATE TABLE IF NOT EXISTS machines (
    id          VARCHAR(64) PRIMARY KEY,
    code        VARCHAR(64) NOT NULL UNIQUE,
    workshop    VARCHAR(120) NOT NULL,
    eq_type     VARCHAR(64) NOT NULL DEFAULT 'pump',
    max_temp_c  DOUBLE PRECISION NOT NULL DEFAULT 85.0,
    max_vibro   DOUBLE PRECISION NOT NULL DEFAULT 7.1,
    
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS telemetry_events (
    id          SERIAL PRIMARY KEY,
    event_id    VARCHAR(64) NOT NULL UNIQUE,
    machine_id  VARCHAR(64) NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
    temp_c      DOUBLE PRECISION NOT NULL,
    vibro       DOUBLE PRECISION NOT NULL,
    rpm         INT NOT NULL,
    power_kw    DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    
    is_overload BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS shift_reports (
    id              VARCHAR(64) PRIMARY KEY,
    machine_id      VARCHAR(64) NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
    status          VARCHAR(32) NOT NULL,
    avg_temp_c      DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    max_temp_c      DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    avg_vibro       DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    max_vibro       DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    overload_count  INT NOT NULL DEFAULT 0,
    risk_level      VARCHAR(32) NOT NULL DEFAULT 'LOW',

    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_telemetry_machine_id ON telemetry_events(machine_id, created_at DESC);
