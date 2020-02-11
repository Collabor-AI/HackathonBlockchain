package services

/*

The structs are arranged in a bottom-up order 

*/

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

type WorldState struct {
	EnsembleMethod string `json:"ensembleMethod,omitempty"`
}

type InitData struct {
	Dataset Dataset `json:"dataset,omitempty"` 
	Objective Objective `json:"objective,omitempty"`
	WorldState WorldState `json:"worldstate,omitempty"`
}

type Transaction struct{
	ID []byte
	Vin []TxInput
	Vout TxOutput
}

type TxInput struct{
	TxID []byte
}

type TxOutput struct{
	Data []byte
}

type BlockData struct {
	Address string `json:"address"`
	Email string `json:"email"`
	TrainPreds string `json:"trainPreds"`
	Preds string `json:"preds"`
	Description string `json:"description,omitempty"`
}

type Wallet struct {
	PrivateKey string `json:"privateKey"`
	PubKey string `json:"pubKey"`
	Address string `json:"address"`
}

type Block struct {
	Timestamp int64
	Data []byte
	PrevBlockHash []byte
	Hash []byte
	Nonce float64
	// Score int
}

type Blockchain struct {
	Tip []byte	
}

type BlockchainIter struct {
	Blocks [][]byte `json:"blocks"`
}


type BlockchainIterator struct {
	currentHash []byte
	// db *badger.DB
}
