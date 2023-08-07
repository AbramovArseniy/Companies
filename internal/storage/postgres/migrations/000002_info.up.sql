CREATE TABLE info(
    address TEXT,
    phone_number VARCHAR(16),
    contant_person VARCHAR(256),
    node_id INT REFERENCES nodes(id) ON DELETE CASCADE
)