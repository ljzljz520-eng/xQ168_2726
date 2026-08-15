package recommend

import (
	"fmt"
	"sort"

	"craftmaterials/internal/domain"
)

type Service struct {
	courses   map[string]domain.Course
	materials []domain.Material
}

func NewService(courses []domain.Course, materials []domain.Material) *Service {
	courseIndex := make(map[string]domain.Course, len(courses))
	for _, course := range courses {
		courseIndex[course.ID] = course
	}
	return &Service{courses: courseIndex, materials: append([]domain.Material(nil), materials...)}
}

func (s *Service) Recommend(learner domain.Learner, limit int) []domain.Recommendation {
	if limit <= 0 {
		return []domain.Recommendation{}
	}
	purchased := stringSet(learner.PurchasedMaterialIDs)
	preferredKinds := kindSet(learner.PreferredKinds)
	courseTags := make(map[string]struct{})
	for _, courseID := range learner.EnrolledCourseIDs {
		course, ok := s.courses[courseID]
		if !ok {
			continue
		}
		for _, tag := range course.Tags {
			courseTags[tag] = struct{}{}
		}
	}

	recommendations := make([]domain.Recommendation, 0, len(s.materials))
	for _, material := range s.materials {
		if _, exists := purchased[material.ID]; exists {
			continue
		}
		score := material.Rating
		reasons := []string{fmt.Sprintf("catalog rating %.1f/5", material.Rating)}

		matches := 0
		for _, tag := range material.Tags {
			if _, ok := courseTags[tag]; ok {
				matches++
			}
		}
		if matches > 0 {
			score += float64(matches) * 2
			reasons = append(reasons, fmt.Sprintf("matches %d enrolled-course topics", matches))
		}
		if material.Difficulty == learner.PreferredDifficulty {
			score += 1.5
			reasons = append(reasons, "matches preferred difficulty")
		}
		if _, ok := preferredKinds[material.Kind]; ok {
			score += 1
			reasons = append(reasons, "matches a preferred material type")
		}
		if rating, ok := learner.Ratings[material.Kind]; ok {
			score += rating / 2
			reasons = append(reasons, fmt.Sprintf("learner rated this material type %.1f/5", rating))
		}

		recommendations = append(recommendations, domain.Recommendation{
			MaterialID:   material.ID,
			MaterialName: material.Name,
			Kind:         material.Kind,
			Score:        score,
			Reasons:      reasons,
		})
	}

	sort.Slice(recommendations, func(i, j int) bool {
		if recommendations[i].Score == recommendations[j].Score {
			return recommendations[i].MaterialID < recommendations[j].MaterialID
		}
		return recommendations[i].Score > recommendations[j].Score
	})
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}
	return recommendations
}

func stringSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func kindSet(values []domain.MaterialKind) map[domain.MaterialKind]struct{} {
	set := make(map[domain.MaterialKind]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}
