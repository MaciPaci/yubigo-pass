CREATE TABLE IF NOT EXISTS passwords
(
    id       Text PRIMARY KEY,
    user_id  TEXT NOT NULL,
    title    TEXT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    url      TEXT,
    nonce    BLOB,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
