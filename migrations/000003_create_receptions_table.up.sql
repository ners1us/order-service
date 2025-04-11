CREATE TABLE receptions
(
    id        UUID PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    pvz_id    UUID      NOT NULL REFERENCES pvzs (id),
    status    TEXT      NOT NULL CHECK (status IN ('in_progress', 'closed'))
);

CREATE INDEX idx_receptions_pvz_id_date_time ON receptions (pvz_id, date_time);
