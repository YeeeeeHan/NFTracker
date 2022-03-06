CREATE TABLE IF NOT EXISTS homes (
    id SERIAL PRIMARY KEY,
    price BIGINT NOT NULL,
    description VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    agent_id BIGINT REFERENCES agents(id)
);