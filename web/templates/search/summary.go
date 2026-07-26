package search

type SearchResultItem struct {
	Type     string `json:"type"` // "project" | "task" | "note"
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	URL      string `json:"url"`
	Badge    string `json:"badge"`
	Icon     string `json:"icon"`
}
