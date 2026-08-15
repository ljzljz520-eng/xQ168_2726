package csvinput

import (
	"errors"
	"strings"
	"testing"
)

func TestParseLearners(t *testing.T) {
	csvData := "learner_id,enrolled_courses,purchased_materials,preferred_difficulty,preferred_kinds,ratings\n" +
		"learner-1,course-1;course-2,mat-1,beginner,fabric;toolkit,fabric=4.5;toolkit=3.5\n"

	learners, err := ParseLearners(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("ParseLearners() error = %v", err)
	}
	if len(learners) != 1 {
		t.Fatalf("len(learners) = %d, want 1", len(learners))
	}
	learner := learners[0]
	if learner.ID != "learner-1" || len(learner.EnrolledCourseIDs) != 2 || len(learner.PurchasedMaterialIDs) != 1 {
		t.Fatalf("learner = %#v", learner)
	}
}

func TestParseLearnersRejectsInvalidRating(t *testing.T) {
	csvData := "learner_id,enrolled_courses,purchased_materials,preferred_difficulty,preferred_kinds,ratings\n" +
		"learner-1,course-1,,beginner,fabric,fabric=5.5\n"

	_, err := ParseLearners(strings.NewReader(csvData))
	if !errors.Is(err, ErrInvalidCSV) {
		t.Fatalf("ParseLearners() error = %v, want %v", err, ErrInvalidCSV)
	}
}
