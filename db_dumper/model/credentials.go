package model
type Credentials struct {
	Dumps []Dump `json:"dumps"`
}

type Dump struct {
	Filename    string `json:"filename"`
	Deleted     bool `json:"deleted"`
	Size        uint64 `json:"size"`
	DownloadURL string `json:"download_url"`
	CreatedAt   string `json:"created_at"`
	DumpID      int `json:"dump_id"`
	ShowURL     string `json:"show_url"`
}