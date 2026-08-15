package fixture

import "craftmaterials/internal/domain"

func Courses() []domain.Course {
	return []domain.Course{
		{ID: "course-sashiko", Name: "Sashiko Coaster", Difficulty: domain.DifficultyBeginner, Tags: []string{"cotton", "hand-sewing", "sashiko"}},
		{ID: "course-tote", Name: "Lined Tote Bag", Difficulty: domain.DifficultyIntermediate, Tags: []string{"canvas", "machine-sewing", "bag"}},
		{ID: "course-doll", Name: "Wool Felt Doll", Difficulty: domain.DifficultyIntermediate, Tags: []string{"felt", "hand-sewing", "doll"}},
	}
}

func Materials() []domain.Material {
	return []domain.Material{
		{ID: "mat-cotton", Name: "Indigo Cotton Bundle", Kind: domain.KindFabric, Difficulty: domain.DifficultyBeginner, Tags: []string{"cotton", "sashiko"}, Rating: 4.8},
		{ID: "mat-needle", Name: "Hand Sewing Toolkit", Kind: domain.KindToolkit, Difficulty: domain.DifficultyBeginner, Tags: []string{"hand-sewing", "sashiko"}, Rating: 4.7},
		{ID: "mat-tote", Name: "Pre-cut Tote Kit", Kind: domain.KindSemiMade, Difficulty: domain.DifficultyIntermediate, Tags: []string{"canvas", "bag"}, Rating: 4.6},
		{ID: "mat-felt", Name: "Wool Felt Palette", Kind: domain.KindFabric, Difficulty: domain.DifficultyIntermediate, Tags: []string{"felt", "doll"}, Rating: 4.5},
		{ID: "mat-cutter", Name: "Rotary Cutter Set", Kind: domain.KindToolkit, Difficulty: domain.DifficultyAdvanced, Tags: []string{"machine-sewing", "canvas"}, Rating: 4.2},
		{ID: "mat-doll", Name: "Stuffed Doll Blank", Kind: domain.KindSemiMade, Difficulty: domain.DifficultyIntermediate, Tags: []string{"doll", "hand-sewing"}, Rating: 4.4},
	}
}

func ExpectedMaterials() map[string][]string {
	return map[string][]string{
		"learner-aya": {"mat-cotton", "mat-needle", "mat-tote"},
		"learner-bo":  {"mat-felt", "mat-doll", "mat-needle"},
	}
}

const LearnerCSV = `learner_id,enrolled_courses,purchased_materials,preferred_difficulty,preferred_kinds,ratings
learner-aya,course-sashiko;course-tote,,beginner,fabric;toolkit,fabric=4.9;toolkit=4.6;semi_finished=4.1
learner-bo,course-doll,mat-cotton,intermediate,semi_finished;fabric,fabric=4.7;toolkit=3.8;semi_finished=4.8
`
