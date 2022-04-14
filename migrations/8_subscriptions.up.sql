-- CREATE TYPE conditions AS ENUM ('above', 'below', 'periodic');

CREATE TABLE subscriptions
(
    id           SERIAL PRIMARY KEY,
    username     VARCHAR REFERENCES users (username),
    slug_name    VARCHAR REFERENCES slugs (slug_name),
--     condition    conditions NOT NULL,
    target_price FLOAT      NOT NULL
);

-- INSERT INTO subscriptions
-- VALUES (1, 'yeeeeehan', 'bayc', 'above', '100');