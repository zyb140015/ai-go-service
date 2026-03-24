package service

import "ai-go-service/internal/domain"

const (
	// DefaultNotePage is the first page returned by list endpoints.
	DefaultNotePage int32 = 1
	// DefaultNotePageSize is the default number of notes returned per request.
	DefaultNotePageSize int32 = 20
	// MaxNotePageSize caps page size to prevent unbounded queries.
	MaxNotePageSize int32 = 100
	// NoteSortFieldCreatedAt sorts notes by creation time.
	NoteSortFieldCreatedAt = "created_at"
	// NoteSortFieldTitle sorts notes by title.
	NoteSortFieldTitle = "title"
	// NoteSortDirectionAsc sorts ascending.
	NoteSortDirectionAsc = "asc"
	// NoteSortDirectionDesc sorts descending.
	NoteSortDirectionDesc = "desc"
)

// NoteListOptions describes supported list-note query parameters.
type NoteListOptions struct {
	Page          int32
	PageSize      int32
	Query         string
	SortField     string
	SortDirection string
}

// NoteListResult contains the current page of notes and its paging metadata.
type NoteListResult struct {
	Items         []domain.Note
	Total         int64
	Page          int32
	PageSize      int32
	Query         string
	SortField     string
	SortDirection string
}
