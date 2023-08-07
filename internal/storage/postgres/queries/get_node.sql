-- name: GetOneNode :one
SELECT nodes.*, info.address, info.phone_number, info.contant_person FROM nodes LEFT JOIN info ON nodes.id = info.node_id;