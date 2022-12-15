package models

import "time"

type Assets struct {
	Assets []*Asset `json:"assets,omitempty"`
}

type SingleAsset struct {
	Asset *Asset `json:"asset,omitempty"`
}

type Asset struct {
	ID           int64       `json:"id,omitempty" format:"int64"`
	DisplayID    int64       `json:"display_id,omitempty" format:"int64"`
	Name         string      `json:"name,omitempty"`
	Description  interface{} `json:"description,omitempty"`
	AssetTypeID  int64       `json:"asset_type_id,omitempty" format:"int64"`
	Impact       string      `json:"impact,omitempty"`
	AuthorType   string      `json:"author_type,omitempty"`
	UsageType    string      `json:"usage_type,omitempty"`
	AssetTag     string      `json:"asset_tag,omitempty"`
	UserID       interface{} `json:"user_id,omitempty"`
	DepartmentID interface{} `json:"department_id,omitempty"`
	LocationID   interface{} `json:"location_id,omitempty"`
	AgentID      interface{} `json:"agent_id,omitempty"`
	GroupID      interface{} `json:"group_id,omitempty"`
	AssignedOn   interface{} `json:"assigned_on,omitempty"`
	CreatedAt    *time.Time  `json:"created_at,omitempty"`
	UpdatedAt    *time.Time  `json:"updated_at,omitempty"`
}
