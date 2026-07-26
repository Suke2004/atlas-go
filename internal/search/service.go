package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/notes"
	"github.com/Suke2004/atlas-go/internal/projects"
	"github.com/Suke2004/atlas-go/internal/tasks"
	searchtemplates "github.com/Suke2004/atlas-go/web/templates/search"
)

type SearchResultItem = searchtemplates.SearchResultItem

type Service interface {
	Search(ctx context.Context, userID int64, query string) ([]SearchResultItem, error)
}

type service struct {
	database    *db.DB
	projectsSvc projects.Service
	tasksSvc    tasks.Service
	notesSvc    notes.Service
}

func NewService(database *db.DB, projSvc projects.Service, tasksSvc tasks.Service, notesSvc notes.Service) Service {
	return &service{
		database:    database,
		projectsSvc: projSvc,
		tasksSvc:    tasksSvc,
		notesSvc:    notesSvc,
	}
}

func (s *service) Search(ctx context.Context, userID int64, query string) ([]SearchResultItem, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, nil
	}

	var results []SearchResultItem

	// 1. Search Projects
	if s.projectsSvc != nil {
		projList, err := s.projectsSvc.ListProjects(ctx, userID, "all", "all", q)
		if err == nil {
			for _, p := range projList {
				results = append(results, SearchResultItem{
					Type:     "project",
					ID:       p.ID,
					Title:    p.Name,
					Subtitle: p.Description,
					URL:      fmt.Sprintf("/projects/%d", p.ID),
					Badge:    "Project",
					Icon:     "folder",
				})
			}
		}
	}

	// 2. Search Tasks
	if s.tasksSvc != nil {
		taskList, err := s.tasksSvc.ListTasks(ctx, userID, "all", "all", "all", 0, q)
		if err == nil {
			for _, t := range taskList {
				results = append(results, SearchResultItem{
					Type:     "task",
					ID:       t.Task.ID,
					Title:    t.Task.Title,
					Subtitle: t.Task.Description,
					URL:      "/tasks",
					Badge:    strings.ToUpper(t.Task.Priority),
					Icon:     "check-square",
				})
			}
		}
	}

	// 3. Search Notes
	if s.notesSvc != nil {
		notesList, err := s.notesSvc.ListNotes(ctx, userID, "all", q)
		if err == nil {
			for _, n := range notesList {
				results = append(results, SearchResultItem{
					Type:     "note",
					ID:       n.Note.ID,
					Title:    n.Note.Title,
					Subtitle: n.Note.Content,
					URL:      fmt.Sprintf("/notes/%d", n.Note.ID),
					Badge:    "Note",
					Icon:     "book-open",
				})
			}
		}
	}

	return results, nil
}
