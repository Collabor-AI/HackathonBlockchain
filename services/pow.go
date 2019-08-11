package services

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBits = 16
const maxNonce = math.MaxInt64


type ProofOfWork struct {
	block *Block
	//score float64
	target *big.Int 	
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}


func IntToHex(Int int64) []byte{
	hex := fmt.Sprintf("%x",Int)
	return []byte(hex)
}

//pow is to find a number smaller than the difficulty

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte){
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf("Mining block containing \"%+v\n", pow.block.Data)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:]) //hashInt of the content of the hashed block

		if hashInt.Cmp(pow.target) == -1 { //loop until find a nonce, which inconjuction with teh block data hashes to target
			break
		} else {
			nonce++ 
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}


func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := (hashInt.Cmp(pow.target) == -1)

	return isValid
}