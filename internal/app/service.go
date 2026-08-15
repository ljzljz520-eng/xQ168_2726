package app

import (
	"io"

	"craftmaterials/internal/csvinput"
	"craftmaterials/internal/domain"
	"craftmaterials/internal/recommend"
	"craftmaterials/internal/report"
)

type LearnerResult struct {
	LearnerID       string                  `json:"learner_id"`
	Recommendations []domain.Recommendation `json:"recommendations"`
}

type Response struct {
	Results    []LearnerResult   `json:"results"`
	Evaluation report.Evaluation `json:"evaluation"`
}

type Service struct {
	recommender *recommend.Service
	expected    map[string][]string
	materialIDs []string
}

func NewService(courses []domain.Course, materials []domain.Material, expected map[string][]string) *Service {
	materialIDs := make([]string, 0, len(materials))
	for _, material := range materials {
		materialIDs = append(materialIDs, material.ID)
	}
	return &Service{
		recommender: recommend.NewService(courses, materials),
		expected:    expected,
		materialIDs: materialIDs,
	}
}

func (s *Service) FromCSV(r io.Reader, limit int) (Response, error) {
	learners, err := csvinput.ParseLearners(r)
	if err != nil {
		return Response{}, err
	}

	response := Response{Results: make([]LearnerResult, 0, len(learners))}
	byLearner := make(map[string][]domain.Recommendation, len(learners))
	for _, learner := range learners {
		recommendations := s.recommender.Recommend(learner, limit)
		response.Results = append(response.Results, LearnerResult{
			LearnerID:       learner.ID,
			Recommendations: recommendations,
		})
		byLearner[learner.ID] = recommendations
	}
	response.Evaluation = report.Evaluate(byLearner, expectedForLearners(s.expected, learners), s.materialIDs)
	return response, nil
}

func expectedForLearners(expected map[string][]string, learners []domain.Learner) map[string][]string {
	selected := make(map[string][]string, len(learners))
	for _, learner := range learners {
		selected[learner.ID] = append([]string(nil), expected[learner.ID]...)
	}
	return selected
}
