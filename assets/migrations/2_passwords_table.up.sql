CREATE TABLE IF NOT EXISTS passwords
(
    user_id  TEXT NOT NULL,
    title    TEXT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    url      TEXT,
    nonce    BLOB,
    PRIMARY KEY (user_id, title, username),
    FOREIGN KEY (user_id) REFERENCES users (id)
);
