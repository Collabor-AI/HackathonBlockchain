package services

type InitData struct {
	Dataset Dataset `json:"dataset,omitempty"` 
	Objective Objective `json:"objective,omitempty"`
	
}

type Dataset struct {
	Method string `json:"method,omitempty"` // Kaggle or URL
	Name string `json:"name,omitempty"` // Name of resource for Kaggle
	Description string `json:"description,omitempty"` //Description of Data, Data Specification
}

type Objective struct {
	Baseline float64 `json:"baseline,omitempty"`
	Scoring string `json:"scoring,omitempty"`
}

type ModelData struct {
	Model []byte `json:"model",omitempty`
	Preds []byte `json:"preds",omitempty`
}

