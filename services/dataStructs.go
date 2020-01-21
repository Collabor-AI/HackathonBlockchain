package services

type InitData struct {
	Dataset Dataset `json:"dataset,omitempty"` 
	Objective Objective `json:"objective,omitempty"`
	WorldState WorldState `json:"worldstate,omitempty"`
}

type Dataset struct {
	TrainURL string `json:"trainURL,omitempty"`
	TestURL string `json:"testURL,omitempty"`
	Name string `json:"name,omitempty"` // Name of resource
	Description string `json:"description,omitempty"` //Description of Data, Data Specification
}

type Objective struct {
	Baseline float64 `json:"baseline,omitempty"` //reject if scores below this
	Scoring string `json:"scoring,omitempty"` //scoring method
}

type BlockData struct {
	Address string `json:"address"`
	Name string  `json:"name"`
	Email string `json:"email"`
	Preds []byte `json:"preds"`
	LinkToCode string `json:"linkToCode,omitempty"`
	Description string `json:"description,omitempty"`
	PrivKey string `json:"privateKey"`
}

type WorldState struct {
	EnsembleMethod string `json:"ensembleMethod,omitempty"`
}