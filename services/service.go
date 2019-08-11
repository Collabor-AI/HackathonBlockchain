package services


import (
	"context"
	// "errors"
	"github.com/dgraph-io/badger"
	// "github.com/go-kit/kit/log"
	"log"
	
	// "strconv"
)


type Service interface {
	NewBlockchain(ctx context.Context, startingData InitData) (*Blockchain, error)
	// AddBlock(ctx context.Context, data modelData, db *badger.DB) (error)
	// PrintBlockchain(ctx context.Context, db *badger.DB) (error)
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

// func (s basicService) PrintBlockchain(ctx context.Context, db *badger.DB) error{
// 	bc := s.NewBlockchain()
// 	bci := bc.Iterator()
// 	for {
// 		block := bci.Next()
// 		if len(block.PrevBlockHash) == 0 {
// 			fmt.Printf("%+v",block)
// 			fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
// 			fmt.Printf("Data: %+v\n", block.Data)
// 			fmt.Printf("Hash: %x\n", block.Hash)
// 			pow := NewProofOfWork(block)
// 			fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
// 			break
// 		}

		

// 		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
// 		fmt.Printf("Data: %s\n", block.Data)
// 		fmt.Printf("Hash: %x\n", block.Hash)
// 		pow := NewProofOfWork(block)
// 		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
// 		fmt.Println()

// 	}
// }

func (s basicService) NewBlockchain(ctx context.Context, startingData InitData) (*Blockchain, error) {

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
			genesis := NewGenesisBlock()
			_ = txn.Set(genesis.Hash, genesis.Serialize())
			_ = txn.Set([]byte("last"), genesis.Hash)
			tip = genesis.Hash

			return nil
		})
	}


	bc := Blockchain{tip}
	if err != nil {
		log.Print("failed to initiate blockchain")
	}

	return &bc, nil
}