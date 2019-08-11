package services

type InitData struct {
	Dataset []byte `json:"dataset"`
	Baseline float64 `json:"baseline"`
}

type ModelData struct {
	model []byte `json:"model",omitempty`
	preds []byte `json:"preds",omitempty`
}

