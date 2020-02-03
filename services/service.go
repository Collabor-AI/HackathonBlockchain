package services


import (
	"context"
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	// "github.com/go-kit/kit/log"
	"log"
	// "strconv"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)


type Service interface {
	NewBlockchain(ctx context.Context, startingData InitData) (Blockchain, error)
	AddBlock(ctx context.Context, bd BlockData, poml float64) (error)
	PrintBlockchain(ctx context.Context) (*BlockchainIter, error)
	GenerateAddress(ctx context.Context)  (*Wallet, error)
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
	//For now our service only needs the blockchain - which is in persistence, badger DB
	db *badger.DB
}

func (s basicService) NewBlockchain(ctx context.Context, startingData InitData) (Blockchain, error) {
	//Check if there is an existing blockchain
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
	//return the tip of the blockchain
	return bc, nil
}


func (s basicService) AddBlock(ctx context.Context, bd BlockData, poml float64) error{
	//Check if there's an existing blockchain
	bc,_ := s.NewBlockchain(ctx, InitData{})
	var lastHash []byte

	// retrieve end of blockchain from badger
	_ = s.db.View(func(txn *badger.Txn) error {
		item, _ := txn.Get([]byte("last"))

		lastHash, _ = item.ValueCopy(nil)
		return nil
	})

	//create a new block
	bdBytes,_ := json.Marshal(bd)
	newBlock := NewBlock(bdBytes, poml, lastHash)

	
	_ = s.db.Update(func(txn *badger.Txn) error {
		_ = txn.Set(newBlock.Hash, newBlock.Serialize())
		_ = txn.Set([]byte("last"), newBlock.Hash)
		bc.Tip = newBlock.Hash
		return nil
	})
	return nil
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
			break
		}

		

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()

	}
	return blocks, nil
}

func (s basicService) GenerateAddress(ctx context.Context) (*Wallet, error){
	privateKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
    publicKey := &privateKey.PublicKey
    encPriv, encPub := encode(privateKey, publicKey)



    pubKeyBytes := elliptic.Marshal(elliptic.P384(), publicKey.X, publicKey.Y)
    address := "0x" + hex.EncodeToString(Keccak256(pubKeyBytes[1:])[12:]) 



	// pubHash := PublicKeyHash(public)

	// versionedHash := append([]byte{version}, pubHash...)
	// checksum := Checksum(versionedHash)

	// fullHash := append(versionedHash, checksum...)
	// address := Base58Encode(fullHash)


	// privateKeyBytes, publicKeyBytes := encode(private, public)
	wallet := Wallet{encPriv, encPub, address}
	// fmt.Printf("Service: %+v\n",wallet)

	return &wallet, nil
}






