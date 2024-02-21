CREATE TABLE IF NOT EXISTS users
(
    id       TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT        NOT NULL,
    salt     TEXT        NOT NULL
);
