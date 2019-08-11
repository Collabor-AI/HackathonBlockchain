package main


type initData struct {
	dataset []byte `json:"dataset"`
	baseline Int `json:"baseline"`
}

type modelData struct {
	model []byte `json:"model",omitempty`
	preds []byte `json:"preds",omitempty`
}

