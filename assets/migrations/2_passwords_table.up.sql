CREATE TABLE IF NOT EXISTS passwords
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT,
    site     TEXT,
    password TEXT,
    UNIQUE(username, site)
);
