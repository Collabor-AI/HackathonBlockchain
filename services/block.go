package services
import (
	// 2
	"crypto/sha256"
	"encoding/json"
	"bytes"
	"encoding/binary"
	"math"
	// "fmt"
	"github.com/dgraph-io/badger"
	"log"
	"strconv"
	"fmt"
	"time"
	
)

func float64ToByte(f float64) []byte {
	//Util to convert float64 score to hash
   var buf [8]byte
   binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
   return buf[:]
}
	
func NewBlock(data []byte, poml float64, prevBlockHash []byte) *Block {
	//A Block Hash is a Sha256 hash of the previous hash, the block data, the nonce
	block := &Block{
		Timestamp:time.Now().Unix(), 
		Data: data,
		PrevBlockHash:prevBlockHash,
		Hash:[]byte{},
		Nonce: poml,
	}

	hashData := bytes.Join(
		[][]byte{
			block.PrevBlockHash,
			block.Data,
			[]byte(strconv.FormatInt(block.Timestamp, 16)),
			float64ToByte(poml),
		},
		[]byte{},
	)

	hash := sha256.Sum256(hashData)
	block.Hash = hash[:]
	return block
}

func NewGenesisBlock(startingData InitData) *Block {
	fmt.Printf("startingData is s %+v",startingData )
	dataBytes,_ := json.Marshal(startingData)
	return NewBlock(dataBytes, startingData.Objective.Baseline,[]byte{})
}

func (b Block) Serialize() []byte{
	blockBytes, _ := json.Marshal(b)
	return blockBytes
}

func DeserializeBlock(d []byte) *Block{
	var block Block
	_ = json.Unmarshal(d,&block)

	return &block
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.Tip}
	return bci
}


func (i *BlockchainIterator) Next(db *badger.DB) *Block {
	var block *Block

	_ = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(i.currentHash)
		if err != nil {
			log.Print("Failed to find", err)
		}
		encodedBlock, _ := item.ValueCopy(nil)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	i.currentHash = block.PrevBlockHash

	return block
}

// func (bc *Blockchain) AddBlock(md modelData){
// 	var lastHash []byte

// 	_ = bc.db.View(func(txn *badger.Txn) error {
// 		item, _ := txn.Get([]byte("last"))

// 		lastHash, _ = item.ValueCopy(nil)
// 		return nil
// 	})

// 	newBlock := NewBlock(md,lastHash)

// 	_ = bc.db.Update(func(txn *badger.Txn) error {
// 		_ = txn.Set(newBlock.Hash, newBlock.Serialize())
// 		_ = txn.Set([]byte("last"), newBlock.Hash)
// 		bc.tip = newBlock.Hash
// 		return nil
// 	})
	
// }






