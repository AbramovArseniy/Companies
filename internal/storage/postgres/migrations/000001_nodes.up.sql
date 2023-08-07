CREATE TABLE nodes(
    id SERIAL UNIQUE,
    name VARCHAR NOT NULL,
    parent_id INT REFERENCES nodes(id) ON DELETE CASCADE
)