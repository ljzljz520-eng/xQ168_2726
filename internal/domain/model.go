package domain

type Difficulty string

const (
	DifficultyBeginner     Difficulty = "beginner"
	DifficultyIntermediate Difficulty = "intermediate"
	DifficultyAdvanced     Difficulty = "advanced"
)

type MaterialKind string

const (
	KindFabric   MaterialKind = "fabric"
	KindToolkit  MaterialKind = "toolkit"
	KindSemiMade MaterialKind = "semi_finished"
)

type Course struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Difficulty Difficulty `json:"difficulty"`
	Tags       []string   `json:"tags"`
}

type Material struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Kind       MaterialKind `json:"kind"`
	Difficulty Difficulty   `json:"difficulty"`
	Tags       []string     `json:"tags"`
	Rating     float64      `json:"rating"`
}

type Learner struct {
	ID                   string                   `json:"id"`
	EnrolledCourseIDs    []string                 `json:"enrolled_course_ids"`
	PurchasedMaterialIDs []string                 `json:"purchased_material_ids"`
	PreferredDifficulty  Difficulty               `json:"preferred_difficulty"`
	PreferredKinds       []MaterialKind           `json:"preferred_kinds"`
	Ratings              map[MaterialKind]float64 `json:"ratings"`
}

type Recommendation struct {
	MaterialID   string       `json:"material_id"`
	MaterialName string       `json:"material_name"`
	Kind         MaterialKind `json:"kind"`
	Score        float64      `json:"score"`
	Reasons      []string     `json:"reasons"`
}
