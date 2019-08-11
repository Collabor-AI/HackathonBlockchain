package services
import (
	// 2
	// "crypto/sha256"
	"encoding/json"
	// "fmt"
	// "github.com/dgraph-io/badger"
	// "log"
	// "strconv"
	"time"
	
)


type Block struct {
	Timestamp int64
	Data []byte
	PrevBlockHash []byte
	Hash []byte
	Nonce int
}

type Blockchain struct {
	tip []byte	
}
	
func NewBlock(data []byte, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:time.Now().Unix(), 
		Data: data,
		PrevBlockHash:prevBlockHash,
		Hash:[]byte{},
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
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






func NewGenesisBlock() *Block {
	data := InitData{[]byte("test"), 0}
	dataBytes,_ := json.Marshal(data)
	return NewBlock(dataBytes, []byte{})
}

func (b *Block) Serialize() []byte{
	blockBytes, _ := json.Marshal(b)
	return blockBytes
}

func DeserializeBlock(d []byte) *Block{
	var block Block
	_ = json.Unmarshal(d,&block)

	return &block
}



// type BlockchainIterator struct {
// 	currentHash []byte
// 	// db *badger.DB
// }

// func (bc *Blockchain) Iterator() *BlockchainIterator {
// 	bci := &BlockchainIterator{bc.tip, bc.db}
// 	return bci
// }

// func (i *BlockchainIterator) Next() *Block {
// 	var block *Block

// 	_ = i.db.View(func(txn *badger.Txn) error {
// 		item, err := txn.Get(i.currentHash)
// 		if err != nil {
// 			log.Print("Failed to find", err)
// 		}
// 		encodedBlock, _ := item.ValueCopy(nil)
// 		block = DeserializeBlock(encodedBlock)

// 		return nil
// 	})

// 	i.currentHash = block.PrevBlockHash

// 	return block
// }

