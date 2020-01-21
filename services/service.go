package services


import (
	"context"
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	// "golang.org/x/crypto/ripemd160" 
	// "errors"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	// "github.com/go-kit/kit/log"
	"log"
	// "strconv"
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
	db *badger.DB
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

	//insert into badger
	_ = s.db.Update(func(txn *badger.Txn) error {
		_ = txn.Set(newBlock.Hash, newBlock.Serialize())
		_ = txn.Set([]byte("last"), newBlock.Hash)
		bc.Tip = newBlock.Hash
		return nil
	})
	return nil
}



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
			break
		}

		

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
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

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

// func PublicKeyHash(pubKey []byte) []byte {
// 	pubHash := sha256.Sum256(pubKey)

// 	hasher := ripemd160.New()
// 	_, err := hasher.Write(pubHash[:])
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	publicRipMD := hasher.Sum(nil)

// 	return publicRipMD
// }

// func Checksum(payload []byte) []byte {
// 	firstHash := sha256.Sum256(payload)
// 	secondHash := sha256.Sum256(firstHash[:])

// 	return secondHash[:checksumLength]
// }

// func ValidateAddress(address string) bool {
// 	pubKeyHash := Base58Decode([]byte(address))
// 	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
// 	version := pubKeyHash[0]
// 	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
// 	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

// 	return bytes.Compare(actualChecksum, targetChecksum) == 0
// }

func (s basicService) GenerateAddress(ctx context.Context) (*Wallet, error){
	private, public := NewKeyPair()	
	return &Wallet{private, public}, nil
}

type Wallet struct {
	privKey ecdsa.PrivateKey
	pubkey []byte
}



