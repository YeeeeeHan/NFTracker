CREATE TABLE IF NOT EXISTS chats
(
    id          BIGINT PRIMARY KEY,
    title       VARCHAR,
    description VARCHAR,
    invitelink  VARCHAR
);

INSERT INTO chats (id, title, description, invitelink)
VALUES (123, 'grouptitle', 'title', 'invitelink');