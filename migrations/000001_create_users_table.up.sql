CREATE TABLE users
(
    id       UUID PRIMARY KEY,
    email    TEXT UNIQUE NOT NULL,
    password TEXT        NOT NULL,
    role     TEXT        NOT NULL CHECK (role IN ('employee', 'moderator'))
);
