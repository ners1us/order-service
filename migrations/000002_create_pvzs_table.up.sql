CREATE TABLE pvzs
(
    id                UUID PRIMARY KEY,
    registration_date TIMESTAMP NOT NULL,
    city              TEXT      NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань'))
);
