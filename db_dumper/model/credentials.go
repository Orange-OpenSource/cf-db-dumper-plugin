package model
type Credentials struct {
	Dumps []Dump `json:"dumps"`
}

type Dump struct {
	CreatedAt   string `json:"created_at"`
	DownloadURL string `json:"download_url"`
	DumpID      int    `json:"dump_id"`
	Filename    string `json:"filename"`
	ShowURL     string `json:"show_url"`
}