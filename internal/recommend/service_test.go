package recommend

import (
	"testing"

	"craftmaterials/internal/csvinput"
	"craftmaterials/internal/fixture"
	"strings"
)

func TestRecommendUsesEnrollmentPreferenceAndPurchaseHistory(t *testing.T) {
	learners, err := csvinput.ParseLearners(strings.NewReader(fixture.LearnerCSV))
	if err != nil {
		t.Fatalf("ParseLearners() error = %v", err)
	}
	service := NewService(fixture.Courses(), fixture.Materials())

	got := service.Recommend(learners[0], 3)
	if len(got) != 3 {
		t.Fatalf("len(recommendations) = %d, want 3", len(got))
	}
	if got[0].MaterialID != "mat-cotton" || got[1].MaterialID != "mat-needle" {
		t.Fatalf("top recommendations = %s, %s", got[0].MaterialID, got[1].MaterialID)
	}
	if len(got[0].Reasons) < 4 {
		t.Fatalf("reasons = %v", got[0].Reasons)
	}

	second := service.Recommend(learners[1], 6)
	for _, recommendation := range second {
		if recommendation.MaterialID == "mat-cotton" {
			t.Fatalf("purchased material returned: %s", recommendation.MaterialID)
		}
	}
}
