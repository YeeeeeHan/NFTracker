CREATE TYPE conditions AS ENUM ('above', 'below', 'periodic');

CREATE TABLE subscriptions
(
    id           SERIAL PRIMARY KEY,
    user_id      BIGINT REFERENCES users (id),
    slug_id      BIGINT REFERENCES slugs (id),
    condition    conditions NOT NULL,
    target_price FLOAT      NOT NULL
);

INSERT INTO subscriptions
VALUES (1, 1, 1, 'above', '100');