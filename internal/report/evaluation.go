package report

import "craftmaterials/internal/domain"

type ConfusionMatrix struct {
	TruePositive  int `json:"true_positive"`
	FalsePositive int `json:"false_positive"`
	TrueNegative  int `json:"true_negative"`
	FalseNegative int `json:"false_negative"`
}

type Evaluation struct {
	Matrix    ConfusionMatrix `json:"confusion_matrix"`
	Precision float64         `json:"precision"`
	Recall    float64         `json:"recall"`
	Accuracy  float64         `json:"accuracy"`
}

func Evaluate(results map[string][]domain.Recommendation, expected map[string][]string, materialIDs []string) Evaluation {
	matrix := ConfusionMatrix{}
	learners := make(map[string]struct{}, len(results)+len(expected))
	for learnerID := range results {
		learners[learnerID] = struct{}{}
	}
	for learnerID := range expected {
		learners[learnerID] = struct{}{}
	}
	for learnerID := range learners {
		predicted := make(map[string]struct{})
		for _, recommendation := range results[learnerID] {
			predicted[recommendation.MaterialID] = struct{}{}
		}
		actual := make(map[string]struct{})
		for _, materialID := range expected[learnerID] {
			actual[materialID] = struct{}{}
		}
		for _, materialID := range materialIDs {
			_, isPredicted := predicted[materialID]
			_, isActual := actual[materialID]
			switch {
			case isPredicted && isActual:
				matrix.TruePositive++
			case isPredicted:
				matrix.FalsePositive++
			case isActual:
				matrix.FalseNegative++
			default:
				matrix.TrueNegative++
			}
		}
	}

	total := matrix.TruePositive + matrix.FalsePositive + matrix.TrueNegative + matrix.FalseNegative
	return Evaluation{
		Matrix:    matrix,
		Precision: ratio(matrix.TruePositive, matrix.TruePositive+matrix.FalsePositive),
		Recall:    ratio(matrix.TruePositive, matrix.TruePositive+matrix.FalseNegative),
		Accuracy:  ratio(matrix.TruePositive+matrix.TrueNegative, total),
	}
}

func ratio(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	return float64(numerator) / float64(denominator)
}
