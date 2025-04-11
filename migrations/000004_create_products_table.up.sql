CREATE TABLE products
(
    id           UUID PRIMARY KEY,
    date_time    TIMESTAMP NOT NULL,
    type         TEXT      NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID      NOT NULL REFERENCES receptions (id)
);

CREATE INDEX idx_products_reception_id ON products (reception_id);
