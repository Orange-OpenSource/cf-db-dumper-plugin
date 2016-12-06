package model

type Credentials struct {
	Dumps        []Dump   `json:"dumps"`
	DatabaseType string   `json:"database_type,omitempty"`
	DatabaseRef  string   `json:"database_ref,omitempty"`
}

type Dump struct {
	Filename    string   `json:"filename"`
	Deleted     bool     `json:"deleted"`
	Size        float64  `json:"size"`
	DownloadURL string   `json:"download_url"`
	CreatedAt   string   `json:"created_at"`
	DumpID      int      `json:"dump_id"`
	ShowURL     string   `json:"show_url"`
	Tags        []string `json:"tags,omitempty"`
}