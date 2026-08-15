package workflow

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusSubmitted Status = "submitted"
	StatusCompleted Status = "completed"
	StatusArchived  Status = "archived"
)

var (
	ErrDocumentNotFound  = errors.New("business document not found")
	ErrInvalidTransition = errors.New("invalid status transition")
)

type Document struct {
	ID        string `json:"id"`
	LearnerID string `json:"learner_id"`
	CourseID  string `json:"course_id"`
	Status    Status `json:"status"`
}

type MemoryStore struct {
	mu        sync.RWMutex
	documents map[string]Document
}

func NewMemoryStore(documents []Document) *MemoryStore {
	store := &MemoryStore{documents: make(map[string]Document, len(documents))}
	for _, document := range documents {
		store.documents[document.ID] = document
	}
	return store
}

func (s *MemoryStore) Get(id string) (Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	document, ok := s.documents[id]
	if !ok {
		return Document{}, fmt.Errorf("%w: %s", ErrDocumentNotFound, id)
	}
	return document, nil
}

func (s *MemoryStore) Save(document Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.documents[document.ID]; !ok {
		return fmt.Errorf("%w: %s", ErrDocumentNotFound, document.ID)
	}
	s.documents[document.ID] = document
	return nil
}

func (s *MemoryStore) ListByStatus(status Status) []Document {
	s.mu.RLock()
	defer s.mu.RUnlock()
	documents := make([]Document, 0)
	for _, document := range s.documents {
		if document.Status == status {
			documents = append(documents, document)
		}
	}
	sort.Slice(documents, func(i, j int) bool { return documents[i].ID < documents[j].ID })
	return documents
}

type StatusService struct {
	store *MemoryStore
}

func NewStatusService(store *MemoryStore) *StatusService {
	return &StatusService{store: store}
}

func (s *StatusService) Transition(id string, target Status) error {
	document, err := s.store.Get(id)
	if err != nil {
		return err
	}
	if !canTransition(document.Status, target) {
		return fmt.Errorf("%w: %s to %s", ErrInvalidTransition, document.Status, target)
	}
	document.Status = target
	return s.store.Save(document)
}

func canTransition(current, target Status) bool {
	switch current {
	case StatusDraft:
		return target == StatusSubmitted
	case StatusSubmitted:
		return target == StatusCompleted
	case StatusCompleted:
		return target == StatusArchived
	default:
		return false
	}
}

func (s *StatusService) CompletedDocuments() []Document {
	return s.store.ListByStatus(StatusCompleted)
}
