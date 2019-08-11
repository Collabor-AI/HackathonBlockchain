package services


import (
	"context"
	// "errors"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	// "github.com/go-kit/kit/log"
	"log"
	
	"strconv"
)


type Service interface {
	NewBlockchain(ctx context.Context, startingData InitData) (Blockchain, error)
	// AddBlock(ctx context.Context, data modelData, db *badger.DB) (error)
	PrintBlockchain(ctx context.Context) (*BlockchainIter, error)
}

func New(db *badger.DB) Service {
	var svc Service
	{
		svc = NewBasicService(db)
	}
	return svc
}	

func NewBasicService(db *badger.DB) Service {
	return basicService{db}
}

type basicService struct{
	db *badger.DB
}


// func (s basicService) AddBlock(ctx context.Context, data modelData, db *badger.DB) error{
// 	bc := s.NewBlockchain()
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
// 	return nil
// }



func (s basicService) NewBlockchain(ctx context.Context, startingData InitData) (Blockchain, error) {

	var tip []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("last"))
		if err != nil {
			log.Print("Initialising New Blockchain", err)
		} else {
			tip,_ = item.ValueCopy(nil)
		}
		return nil
	})

	if len(tip) == 0 {
		err = s.db.Update(func(txn *badger.Txn) error {
			genesis := NewGenesisBlock(startingData)
			_ = txn.Set(genesis.Hash, genesis.Serialize())
			_ = txn.Set([]byte("last"), genesis.Hash)
			tip = genesis.Hash
			log.Print("genesis is %+v",genesis)
			return nil
		})
	}



	bc := Blockchain{tip}
	if err != nil {
		log.Print("failed to initiate blockchain")
	}

	return bc, nil
}


func (s basicService) PrintBlockchain(ctx context.Context) (*BlockchainIter, error){
	bc,_ := s.NewBlockchain(ctx, InitData{})
	blocks := &BlockchainIter{[][]byte{}}
	bci := bc.Iterator()
	for {
		block := bci.Next(s.db)
		blockData,_ := json.Marshal(block)
		blocks.Blocks = append(blocks.Blocks,blockData)
		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("%+v",block)
			fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
			fmt.Printf("Data: %+v\n", block.Data)
			fmt.Printf("Hash: %x\n", block.Hash)
			pow := NewProofOfWork(block)
			fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
			break
		}

		

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

	}
	return blocks, nil
}

type BlockchainIter struct {
	Blocks [][]byte `json:"blocks"`
}


type BlockchainIterator struct {
	currentHash []byte
	// db *badger.DB
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

