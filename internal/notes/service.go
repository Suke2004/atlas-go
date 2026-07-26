package notes

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/projects"
	notetemplates "github.com/Suke2004/atlas-go/web/templates/notes"
)

type NoteInput struct {
	ProjectID int64    `json:"project_id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	IsPinned  bool     `json:"is_pinned"`
	Tags      []string `json:"tags"`
}

type NoteWithDetails = notetemplates.NoteWithDetails
type NotesSummary = notetemplates.NotesSummary

type Service interface {
	CreateNote(ctx context.Context, userID int64, input NoteInput) (db.Note, error)
	GetNote(ctx context.Context, userID, noteID int64) (*NoteWithDetails, error)
	ListNotes(ctx context.Context, userID int64, tagFilter string, searchQuery string) ([]NoteWithDetails, error)
	GetNotesSummary(ctx context.Context, userID int64) (NotesSummary, error)
	UpdateNote(ctx context.Context, userID, noteID int64, input NoteInput) (db.Note, error)
	TogglePin(ctx context.Context, userID, noteID int64) (db.Note, error)
	DeleteNote(ctx context.Context, userID, noteID int64) error
	GetTemplateContent(templateType string) (string, string)
}

type service struct {
	repo         Repository
	projectsRepo projects.Repository
}

func NewService(repo Repository, projRepo projects.Repository) Service {
	return &service{
		repo:         repo,
		projectsRepo: projRepo,
	}
}

var wikiLinkRegex = regexp.MustCompile(`\[\[(.*?)\]\]`)

func (s *service) CreateNote(ctx context.Context, userID int64, input NoteInput) (db.Note, error) {
	if strings.TrimSpace(input.Title) == "" {
		return db.Note{}, fmt.Errorf("note title is required")
	}

	var projectID interface{}
	if input.ProjectID > 0 {
		projectID = input.ProjectID
	}

	arg := db.CreateNoteParams{
		UserID:    userID,
		ProjectID: projectID,
		Title:     strings.TrimSpace(input.Title),
		Content:   input.Content,
		IsPinned:  input.IsPinned,
	}

	n, err := s.repo.CreateNote(ctx, arg)
	if err != nil {
		return db.Note{}, err
	}

	// Index tags
	for _, tag := range input.Tags {
		t := strings.TrimPrefix(strings.TrimSpace(tag), "#")
		if t != "" {
			_ = s.repo.AddNoteTag(ctx, db.AddNoteTagParams{
				NoteID: n.ID,
				Tag:    t,
			})
		}
	}

	// Parse [[Wiki Links]] and link notes
	s.processWikiLinks(ctx, userID, n.ID, input.Content)

	return n, nil
}

func (s *service) GetNote(ctx context.Context, userID, noteID int64) (*NoteWithDetails, error) {
	n, err := s.repo.GetNoteByID(ctx, db.GetNoteByIDParams{ID: noteID, UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("note not found: %w", err)
	}

	var projName string
	if n.ProjectID.Valid && s.projectsRepo != nil {
		proj, err := s.projectsRepo.GetProjectByID(ctx, db.GetProjectByIDParams{ID: n.ProjectID.Int64, UserID: userID})
		if err == nil {
			projName = proj.Name
		}
	}

	tags, _ := s.repo.ListNoteTags(ctx, noteID)
	backlinks, _ := s.repo.GetNoteBacklinks(ctx, db.GetNoteBacklinksParams{TargetNoteID: noteID, UserID: userID})

	words := len(strings.Fields(n.Content))
	readingTime := (words / 200) + 1

	return &NoteWithDetails{
		Note:        n,
		ProjectName: projName,
		Tags:        tags,
		Backlinks:   backlinks,
		WordCount:   words,
		ReadingTime: readingTime,
	}, nil
}

func (s *service) ListNotes(ctx context.Context, userID int64, tagFilter string, searchQuery string) ([]NoteWithDetails, error) {
	allNotes, err := s.repo.ListNotes(ctx, userID)
	if err != nil {
		return nil, err
	}

	var projectsMap map[int64]string
	if s.projectsRepo != nil {
		projList, err := s.projectsRepo.ListProjects(ctx, userID)
		if err == nil {
			projectsMap = make(map[int64]string)
			for _, p := range projList {
				projectsMap[p.ID] = p.Name
			}
		}
	}

	var filtered []NoteWithDetails
	searchLower := strings.ToLower(strings.TrimSpace(searchQuery))
	tagLower := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(tagFilter), "#"))

	for _, n := range allNotes {
		tags, _ := s.repo.ListNoteTags(ctx, n.ID)

		if tagLower != "" && tagLower != "all" {
			hasTag := false
			for _, t := range tags {
				if strings.EqualFold(t, tagLower) {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		if searchLower != "" {
			titleMatch := strings.Contains(strings.ToLower(n.Title), searchLower)
			contentMatch := strings.Contains(strings.ToLower(n.Content), searchLower)
			if !titleMatch && !contentMatch {
				continue
			}
		}

		projName := ""
		if n.ProjectID.Valid && projectsMap != nil {
			projName = projectsMap[n.ProjectID.Int64]
		}

		words := len(strings.Fields(n.Content))
		readingTime := (words / 200) + 1

		filtered = append(filtered, NoteWithDetails{
			Note:        n,
			ProjectName: projName,
			Tags:        tags,
			WordCount:   words,
			ReadingTime: readingTime,
		})
	}

	return filtered, nil
}

func (s *service) GetNotesSummary(ctx context.Context, userID int64) (NotesSummary, error) {
	allNotes, err := s.repo.ListNotes(ctx, userID)
	if err != nil {
		return NotesSummary{}, err
	}

	var summary NotesSummary
	summary.TotalNotes = int64(len(allNotes))
	tagMap := make(map[string]bool)

	for _, n := range allNotes {
		if n.IsPinned {
			summary.PinnedNotes++
		}
		tags, _ := s.repo.ListNoteTags(ctx, n.ID)
		for _, t := range tags {
			tagMap[t] = true
		}
	}

	for t := range tagMap {
		summary.AllTags = append(summary.AllTags, t)
	}
	summary.TotalTags = int64(len(summary.AllTags))

	return summary, nil
}

func (s *service) UpdateNote(ctx context.Context, userID, noteID int64, input NoteInput) (db.Note, error) {
	if strings.TrimSpace(input.Title) == "" {
		return db.Note{}, fmt.Errorf("note title is required")
	}

	var projectID interface{}
	if input.ProjectID > 0 {
		projectID = input.ProjectID
	}

	n, err := s.repo.UpdateNote(ctx, db.UpdateNoteParams{
		ProjectID: projectID,
		Title:     strings.TrimSpace(input.Title),
		Content:   input.Content,
		IsPinned:  input.IsPinned,
		ID:        noteID,
		UserID:    userID,
	})
	if err != nil {
		return db.Note{}, err
	}

	for _, tag := range input.Tags {
		t := strings.TrimPrefix(strings.TrimSpace(tag), "#")
		if t != "" {
			_ = s.repo.AddNoteTag(ctx, db.AddNoteTagParams{
				NoteID: n.ID,
				Tag:    t,
			})
		}
	}

	s.processWikiLinks(ctx, userID, n.ID, input.Content)

	return n, nil
}

func (s *service) TogglePin(ctx context.Context, userID, noteID int64) (db.Note, error) {
	n, err := s.repo.GetNoteByID(ctx, db.GetNoteByIDParams{ID: noteID, UserID: userID})
	if err != nil {
		return db.Note{}, err
	}

	var projID interface{}
	if n.ProjectID.Valid {
		projID = n.ProjectID.Int64
	}

	return s.repo.UpdateNote(ctx, db.UpdateNoteParams{
		ProjectID: projID,
		Title:     n.Title,
		Content:   n.Content,
		IsPinned:  !n.IsPinned,
		ID:        noteID,
		UserID:    userID,
	})
}

func (s *service) DeleteNote(ctx context.Context, userID, noteID int64) error {
	return s.repo.DeleteNote(ctx, db.DeleteNoteParams{ID: noteID, UserID: userID})
}

func (s *service) GetTemplateContent(templateType string) (string, string) {
	switch strings.ToLower(templateType) {
	case "adr":
		return "ADR: Architecture Decision Title", `# ADR: Title

## Status
Proposed / Accepted / Superseded

## Context
What is the problem or architectural need?

## Decision
What is the chosen design or pattern?

## Consequences
- **Positive**: High throughput, clean separation
- **Negative**: Increased storage footprint`

	case "meeting":
		return "Meeting Notes — [Topic]", `# Meeting Notes — [Date]

## Attendees
- Alex, Developer

## Agenda
1. Architecture review
2. Deliverable timelines

## Key Decisions
- Adopted SQLite WAL mode for high concurrency.

## Action Items
- [ ] Implement Phase 6 Knowledge Base`

	case "brainstorm":
		return "Brainstorming — [Feature]", `# Brainstorming: Feature Idea

## Problem Statement
What user friction are we solving?

## Solution Concept
High level overview of proposed feature.

## User Flow
1. User clicks action button
2. System processes input instantly`

	default:
		return "New Note", `# New Note

Start writing markdown content here...`
	}
}

func (s *service) processWikiLinks(ctx context.Context, userID, sourceNoteID int64, content string) {
	matches := wikiLinkRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return
	}

	allNotes, err := s.repo.ListNotes(ctx, userID)
	if err != nil {
		return
	}

	notesMap := make(map[string]int64)
	for _, n := range allNotes {
		notesMap[strings.ToLower(strings.TrimSpace(n.Title))] = n.ID
	}

	for _, match := range matches {
		if len(match) > 1 {
			targetTitle := strings.ToLower(strings.TrimSpace(match[1]))
			if targetID, exists := notesMap[targetTitle]; exists && targetID != sourceNoteID {
				_ = s.repo.AddNoteLink(ctx, db.AddNoteLinkParams{
					SourceNoteID: sourceNoteID,
					TargetNoteID: targetID,
				})
			}
		}
	}
}
