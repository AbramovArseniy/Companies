package http

import (
	"time"

	"github.com/AbramovArseniy/Companies/internal/storage/postgres/generated/db"
)

type node struct {
	ID            int32  `json:"id"`
	Name          string `json:"name"`
	ParentID      int32  `json:"parent_id,omitempty"`
	Address       string `json:"address,omitempty"`
	PhoneNumber   string `json:"phone_number,omitempty"`
	ContactPerson string `json:"contact_person,omitempty"`
}

type Tag struct {
	UUID          string     `json:"uuid"`
	Name          string     `json:"name"`
	Value         float64    `json:"value"`
	AlertType     string     `json:"alert_type,omitempty"`
	AlertSeverity string     `json:"alert_severity,omitempty"`
	AlertTime     *time.Time `json:"alert_time,omitempty"`
}

func tagRowsToTags(rows []db.GetNodeTagsRow) []Tag {
	tags := make([]Tag, len(rows))
	for i, row := range rows {
		tags[i].UUID = row.Uuid
		tags[i].Name = row.Name
		tags[i].Value = row.Value
		if row.Uuid_2.Valid {
			tags[i].AlertType = row.Type.String
			tags[i].AlertSeverity = row.Severity.String
			t := row.AlertTime.Time
			tags[i].AlertTime = &t
		}
	}
	return tags
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
