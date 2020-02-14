package services
import (
	// 2
	"context"
	"crypto/sha256"
	"encoding/json"
	"bytes"
	"encoding/binary"
	"encoding/base64"
	"math"
	// "fmt"
	"github.com/dgraph-io/badger"
	firebaseDB "firebase.google.com/go/db"
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
	err := json.Unmarshal(d,&block)
	if err != nil {
		return &block	
	} else {
		fmt.Printf("Err %v", err)
		return nil
	}
	
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.Tip}
	return bci
}


func (i *BlockchainIterator) Next(db *badger.DB, client *firebaseDB.Client, ctx context.Context, name string) *Block {
	var block *Block

	// _ = db.View(func(txn *badger.Txn) error {
	// 	key := append([]byte(name),[]byte("-")...)
	// 	key = append(key,i.currentHash...)
	// 	item, err := txn.Get(key)
	// 	if err != nil {
	// 		log.Print("Failed to find", err)
	// 	}
	// 	encodedBlock, _ := item.ValueCopy(nil)
	// 	block = DeserializeBlock(encodedBlock)

	// 	return nil
	// })


	//get name-lastHash from firebase

	// var encodedBlockRaw []byte

	key := base64.StdEncoding.EncodeToString(i.currentHash)
	fmt.Printf("KEY IS : %v",key)
	if err := client.NewRef(key).Get(ctx, &block); err != nil {
		// encodedBlock, _ := base64.StdEncoding.DecodeString(string(encodedBlockRaw))
	 //  block = DeserializeBlock(encodedBlock)
	  log.Print(err)
	  fmt.Print("BLOCK: Failed to get lastHash from firebase\n")
	}
		
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






