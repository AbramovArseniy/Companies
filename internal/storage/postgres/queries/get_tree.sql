-- name: GetAllTree :many
SELECT nodes.*, info.address, info.phone_number, info.contant_person FROM nodes JOIN info ON nodes.id = info.node_id;
