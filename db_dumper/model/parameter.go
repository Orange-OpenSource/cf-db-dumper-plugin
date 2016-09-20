package model

type Parameter struct {
	Action      string    `json:"action,omitempty"`
	Db          string    `json:"db,omitempty"`
	CreatedAt   string    `json:"created_at,omitempty"`
	CfUserToken string    `json:"cf_user_token,omitempty"`
	Org         string    `json:"org,omitempty"`
	Space       string    `json:"space,omitempty"`
	Metadata    *Metadata `json:"metadata,omitempty"`
}
type Metadata struct {
	Tags []string   `json:"tags,omitempty"`
}
type BindingParameter struct {
	SeeAllDumps bool    `json:"see_all_dumps,omitempty"`
	FindByTags  []string    `json:"find_by_tags,omitempty"`
}