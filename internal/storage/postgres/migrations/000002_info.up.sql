CREATE TABLE info(
    address TEXT,
    phone_number VARCHAR(16),
    contact_person VARCHAR(256),
    node_id INT REFERENCES nodes(id) ON DELETE CASCADE 
)