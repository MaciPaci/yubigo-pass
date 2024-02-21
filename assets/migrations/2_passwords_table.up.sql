CREATE TABLE IF NOT EXISTS passwords
(
    user_id  TEXT,
    site     TEXT,
    password TEXT,
    PRIMARY KEY (user_id, site),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
