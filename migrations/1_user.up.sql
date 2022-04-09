CREATE TABLE IF NOT EXISTS users (
    id  SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL
);

INSERT INTO users
VALUES (1, '@yeeeeehan');