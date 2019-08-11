package services

type InitData struct {
	dataset []byte `json:"dataset"`
	baseline int `json:"baseline"`
}

type ModelData struct {
	model []byte `json:"model",omitempty`
	preds []byte `json:"preds",omitempty`
}

