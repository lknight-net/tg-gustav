CREATE TABLE files
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL ,
    fileext TEXT NOT NULL ,
    short_url_original TEXT,
    short_url_upgrade TEXT,
    user_id   INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);