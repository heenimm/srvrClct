CREATE TABLE IF NOT EXISTS users (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     login TEXT UNIQUE NOT NULL,
                                     password_hash TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS expressions (
                                           id INTEGER PRIMARY KEY AUTOINCREMENT,
                                           user_id INTEGER NOT NULL,
                                           expression TEXT NOT NULL,
                                           result TEXT,
                                           created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);