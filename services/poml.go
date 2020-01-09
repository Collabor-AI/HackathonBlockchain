package services

import (
	// "bytes"
	// "crypto/sha256"
	// "fmt"
	// "math/big"
	// "net/http"
)

const(
	PYTHONURL = "localhost:5000/listen"
)



type ProofOfML struct {
	block *Block
	score float64
}

// func (poml *ProofOfML) Run() (int, []byte){
// 	//this will call our python server and return a score
// 	resp,err := http.Get(PYTHONURL)

// }


// func (pow *ProofOfWork) Validate() bool {
// 	var hashInt big.Int

// 	data := pow.prepareData(pow.block.Nonce)
// 	hash := sha256.Sum256(data)
// 	hashInt.SetBytes(hash[:])
// 	isValid := (hashInt.Cmp(pow.target) == -1)

// 	return isValid
// }