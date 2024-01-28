CREATE TABLE IF NOT EXISTS postcards (
    Name VARCHAR(255) NOT NULL,
    path TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id           INTEGER PRIMARY KEY,
    mailing_time INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS discord_servers (
    id BIGINT PRIMARY KEY,
    mailing_time INTEGER NOT NULL
);