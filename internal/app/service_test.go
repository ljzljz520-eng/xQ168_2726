package app

import (
	"strings"
	"testing"

	"craftmaterials/internal/fixture"
)

func TestFixtureRecommendationReport(t *testing.T) {
	service := NewService(fixture.Courses(), fixture.Materials(), fixture.ExpectedMaterials())
	response, err := service.FromCSV(strings.NewReader(fixture.LearnerCSV), 3)
	if err != nil {
		t.Fatalf("FromCSV() error = %v", err)
	}
	if len(response.Results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(response.Results))
	}
	matrix := response.Evaluation.Matrix
	if matrix.TruePositive != 4 || matrix.FalsePositive != 2 || matrix.TrueNegative != 4 || matrix.FalseNegative != 2 {
		t.Fatalf("confusion matrix = %#v", matrix)
	}
	if response.Evaluation.Precision != 2.0/3.0 || response.Evaluation.Recall != 2.0/3.0 {
		t.Fatalf("evaluation = %#v", response.Evaluation)
	}
}
