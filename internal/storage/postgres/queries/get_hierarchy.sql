-- name: GetHierarchy :many
WITH RECURSIVE r(id, name, parent_id, level) AS
        (SELECT tr.id, tr.name, tr.parent_id, 1
        FROM nodes tl
        LEFT JOIN nodes tr 
        ON tl.id = tr.id
        WHERE tl.id = $1
        UNION ALL
        SELECT t.id, t.name, t.parent_id, level+1
        FROM nodes t, r
        WHERE t.id = r.parent_id )
        SELECT id, name, parent_id, ROW_NUMBER() OVER (ORDER BY level DESC) AS level, info.address, info.phone_number, info.contant_person FROM r AS hierarchy 
        LEFT JOIN info 
        ON hierarchy.id = info.node_id;
        