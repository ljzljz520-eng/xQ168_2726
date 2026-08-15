package workflow

import (
	"errors"
	"testing"
)

func TestMaterialRequestLifecycle(t *testing.T) {
	t.Run("submitted request completes", func(t *testing.T) {
		store := NewMemoryStore([]Document{{ID: "request-open", LearnerID: "learner-aya", CourseID: "course-sashiko", Status: StatusSubmitted}})
		service := NewStatusService(store)

		if err := service.Transition("request-open", StatusCompleted); err != nil {
			t.Fatalf("Transition() error = %v", err)
		}
		document, err := store.Get("request-open")
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if document.Status != StatusCompleted {
			t.Fatalf("status = %s, want %s", document.Status, StatusCompleted)
		}
	})

	t.Run("archived request stays out of completed work", func(t *testing.T) {
		store := NewMemoryStore([]Document{{ID: "request-archived", LearnerID: "learner-bo", CourseID: "course-doll", Status: StatusArchived}})
		service := NewStatusService(store)

		err := service.Transition("request-archived", StatusCompleted)
		if !errors.Is(err, ErrInvalidTransition) {
			t.Errorf("Transition() error = %v, want %v", err, ErrInvalidTransition)
		}
		document, getErr := store.Get("request-archived")
		if getErr != nil {
			t.Fatalf("Get() error = %v", getErr)
		}
		if document.Status != StatusArchived {
			t.Errorf("status = %s, want %s", document.Status, StatusArchived)
		}
		completed := service.CompletedDocuments()
		if len(completed) != 0 {
			t.Errorf("completed documents = %#v, want none", completed)
		}
	})
}
