CREATE TABLE alerts(
    type TEXT NOT NULL,
    uuid TEXT NOT NULL,
    alert_time  TIMESTAMP NOT NULL,
    severity VARCHAR(64) NOT NULL,
    state VARCHAR(16) NOT NULL
)