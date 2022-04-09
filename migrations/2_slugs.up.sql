CREATE TABLE IF NOT EXISTS slugs
(
    id                 SERIAL PRIMARY KEY,
    slug_name               VARCHAR NOT NULL,
    floor_price             FLOAT NOT NULL,
    one_day_average_price   BIGINT NOT NULL,
    seven_day_average_price BIGINT NOT NULL
);

INSERT INTO slugs
VALUES (1, 'boredapeyachtclub', 101.2, 101.3, 97.4);
