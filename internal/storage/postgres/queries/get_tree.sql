-- name: GetAllTree :many
SELECT nodes.*, info.address, info.phone_number, info.contact_person FROM nodes LEFT JOIN info ON nodes.id = info.node_id;
