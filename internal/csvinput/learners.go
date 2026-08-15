package csvinput

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"craftmaterials/internal/domain"
)

var ErrInvalidCSV = errors.New("invalid learner CSV")

func ParseLearners(r io.Reader) ([]domain.Learner, error) {
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCSV, err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("%w: header and at least one learner are required", ErrInvalidCSV)
	}
	want := []string{"learner_id", "enrolled_courses", "purchased_materials", "preferred_difficulty", "preferred_kinds", "ratings"}
	if len(records[0]) != len(want) {
		return nil, fmt.Errorf("%w: unexpected header", ErrInvalidCSV)
	}
	for i := range want {
		if strings.TrimSpace(records[0][i]) != want[i] {
			return nil, fmt.Errorf("%w: unexpected header", ErrInvalidCSV)
		}
	}

	learners := make([]domain.Learner, 0, len(records)-1)
	seen := make(map[string]struct{}, len(records)-1)
	for rowIndex, record := range records[1:] {
		if len(record) != len(want) {
			return nil, fmt.Errorf("%w: row %d has %d columns", ErrInvalidCSV, rowIndex+2, len(record))
		}
		id := strings.TrimSpace(record[0])
		if id == "" {
			return nil, fmt.Errorf("%w: row %d has an empty learner_id", ErrInvalidCSV, rowIndex+2)
		}
		if _, ok := seen[id]; ok {
			return nil, fmt.Errorf("%w: duplicate learner_id %q", ErrInvalidCSV, id)
		}
		seen[id] = struct{}{}

		difficulty, err := parseDifficulty(record[3])
		if err != nil {
			return nil, fmt.Errorf("%w: row %d: %v", ErrInvalidCSV, rowIndex+2, err)
		}
		kinds, err := parseKinds(record[4])
		if err != nil {
			return nil, fmt.Errorf("%w: row %d: %v", ErrInvalidCSV, rowIndex+2, err)
		}
		ratings, err := parseRatings(record[5])
		if err != nil {
			return nil, fmt.Errorf("%w: row %d: %v", ErrInvalidCSV, rowIndex+2, err)
		}

		learners = append(learners, domain.Learner{
			ID:                   id,
			EnrolledCourseIDs:    splitList(record[1]),
			PurchasedMaterialIDs: splitList(record[2]),
			PreferredDifficulty:  difficulty,
			PreferredKinds:       kinds,
			Ratings:              ratings,
		})
	}
	return learners, nil
}

func parseDifficulty(value string) (domain.Difficulty, error) {
	difficulty := domain.Difficulty(strings.TrimSpace(value))
	switch difficulty {
	case domain.DifficultyBeginner, domain.DifficultyIntermediate, domain.DifficultyAdvanced:
		return difficulty, nil
	default:
		return "", fmt.Errorf("unknown difficulty %q", value)
	}
}

func parseKinds(value string) ([]domain.MaterialKind, error) {
	parts := splitList(value)
	if len(parts) == 0 {
		return nil, errors.New("preferred_kinds is empty")
	}
	kinds := make([]domain.MaterialKind, 0, len(parts))
	for _, part := range parts {
		kind := domain.MaterialKind(part)
		if !validKind(kind) {
			return nil, fmt.Errorf("unknown material kind %q", part)
		}
		kinds = append(kinds, kind)
	}
	return kinds, nil
}

func parseRatings(value string) (map[domain.MaterialKind]float64, error) {
	ratings := make(map[domain.MaterialKind]float64)
	for _, pair := range splitList(value) {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rating %q", pair)
		}
		kind := domain.MaterialKind(strings.TrimSpace(parts[0]))
		if !validKind(kind) {
			return nil, fmt.Errorf("unknown rating kind %q", kind)
		}
		rating, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil || rating < 0 || rating > 5 {
			return nil, fmt.Errorf("rating for %s must be between 0 and 5", kind)
		}
		ratings[kind] = rating
	}
	return ratings, nil
}

func validKind(kind domain.MaterialKind) bool {
	return kind == domain.KindFabric || kind == domain.KindToolkit || kind == domain.KindSemiMade
}

func splitList(value string) []string {
	parts := strings.Split(strings.TrimSpace(value), ";")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if item := strings.TrimSpace(part); item != "" {
			result = append(result, item)
		}
	}
	return result
}
