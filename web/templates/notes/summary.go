package notes

import "github.com/Suke2004/atlas-go/internal/db"

type NotesSummary struct {
	TotalNotes  int64    `json:"total_notes"`
	PinnedNotes int64    `json:"pinned_notes"`
	TotalTags   int64    `json:"total_tags"`
	AllTags     []string `json:"all_tags"`
}

type NoteWithDetails struct {
	Note        db.Note   `json:"note"`
	ProjectName string    `json:"project_name"`
	Tags        []string  `json:"tags"`
	Backlinks   []db.Note `json:"backlinks"`
	WordCount   int       `json:"word_count"`
	ReadingTime int       `json:"reading_time"`
}
