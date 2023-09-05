package http

import "github.com/AbramovArseniy/Companies/internal/storage/postgres/db"

type node struct {
	ID            int32  `json:"id"`
	Name          string `json:"name"`
	ParentID      int32  `json:"parent_id,omitempty"`
	Address       string `json:"address,omitempty"`
	PhoneNumber   string `json:"phone_number,omitempty"`
	ContactPerson string `json:"contact_person,omitempty"`
}

func treeRowsToNodes(rows []db.GetAllTreeRow) []node {
	nodes := make([]node, len(rows))
	for i, v := range rows {
		nodes[i].ID = v.ID
		nodes[i].Name = v.Name
		if v.Address.Valid {
			nodes[i].Address = v.Address.String
		}
		if v.ParentID.Valid {
			nodes[i].ParentID = v.ParentID.Int32
		}
		if v.PhoneNumber.Valid {
			nodes[i].PhoneNumber = v.PhoneNumber.String
		}
		if v.ContactPerson.Valid {
			nodes[i].ContactPerson = v.ContactPerson.String
		}
	}
	return nodes
}

func hierarchyRowsToNodes(rows []db.GetHierarchyRow) []node {
	nodes := make([]node, len(rows))
	for i, v := range rows {
		if v.ID.Valid {
			nodes[i].ID = v.ID.Int32
		}
		if v.Name.Valid {
			nodes[i].Name = v.Name.String
		}
		if v.Address.Valid {
			nodes[i].Address = v.Address.String
		}
		if v.ParentID.Valid {
			nodes[i].ParentID = v.ParentID.Int32
		}
		if v.PhoneNumber.Valid {
			nodes[i].PhoneNumber = v.PhoneNumber.String
		}
		if v.ContactPerson.Valid {
			nodes[i].ContactPerson = v.ContactPerson.String
		}
	}
	return nodes
}

func nodeRowToNode(row db.GetOneNodeRow) node {
	var n node
	n.ID = row.ID
	n.Name = row.Name
	if row.Address.Valid {
		n.Address = row.Address.String
	}
	if row.ParentID.Valid {
		n.ParentID = row.ParentID.Int32
	}
	if row.PhoneNumber.Valid {
		n.PhoneNumber = row.PhoneNumber.String
	}
	if row.ContactPerson.Valid {
		n.ContactPerson = row.ContactPerson.String
	}
	return n
}
