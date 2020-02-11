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
	//github.com/google/uuid
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
	AddBlock(ctx context.Context, bd BlockData, poml float64, PubKey string ,PrivKey string, Hash string) (error)
	PrintBlockchain(ctx context.Context, Hash string) (*BlockchainIter, error)
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
	
	key := append([]byte(startingData.Dataset.Name),[]byte("-last")...)
	log.Print("Key: %+v",key)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			log.Print("Initialising New Blockchain", err)
		} else {
			tip,_ = item.ValueCopy(nil)
		}
		return nil
	})
	var hashList [][]byte
	if len(tip) == 0 {
		genesis := NewGenesisBlock(startingData)
		err = s.db.Update(func(txn *badger.Txn) error {		
			key2 := append([]byte(startingData.Dataset.Name),[]byte("-")...)	
			key2 = append(key2,genesis.Hash...)
			_ = txn.Set(key2, genesis.Serialize())
			_ = txn.Set(key, genesis.Hash)		
			_ = txn.Set(genesis.Hash, genesis.Serialize())	
			item2, err2 := txn.Get([]byte("hashList"))			
			if err2 != nil {
				log.Print("Failed here1")
				hashList = append(hashList,[][]byte{genesis.Hash}...)
			} else {
				data,_ := item2.ValueCopy(nil)
				_ = json.Unmarshal(data,&hashList)
				hashList = append(hashList,genesis.Hash)
			}
			
			hashListBytes,_ := json.Marshal(hashList)
			_ = txn.Set([]byte("hashList"), hashListBytes)
			tip = genesis.Hash
			log.Print("genesis is %+v",genesis)
			return nil
		})
	}



	bc := Blockchain{tip}
	if err != nil {
		log.Print("failed to initiate blockchain")
	}
	log.Print("BC %+v",bc)
	
	//return the tip of the blockchain
	return bc, nil
}


func (s basicService) AddBlock(ctx context.Context, bd BlockData, poml float64, PubKey string , PrivKey string, Hash string) error{
	var lastHash []byte
	key := append([]byte(Hash),[]byte("-last")...)
	// retrieve end of blockchain from badger
	_ = s.db.View(func(txn *badger.Txn) error {
		item, _ := txn.Get(key)
		lastHash, _ = item.ValueCopy(nil)
		return nil
	})
	log.Printf("PrivKey %v PubKey %v",PrivKey ,PubKey)
	_,pub := decode(PrivKey,PubKey)

	// log.Print("Fialed Here3")
	log.Print("Failed here1 %v %v",key, string(key))
	//create a new block
    pubKeyBytes := elliptic.Marshal(elliptic.P384(), pub.X, pub.Y)
    address := "0x" + hex.EncodeToString(Keccak256(pubKeyBytes[1:])[12:]) 
	bd.Address = address
	bdBytes,_ := json.Marshal(bd)
	newBlock := NewBlock(bdBytes, poml, lastHash)

	log.Print("Failed here3 %v %v",key, string(key))
	_ = s.db.Update(func(txn *badger.Txn) error {
		key2 := append([]byte(Hash),[]byte("-")...)
		key2 = append(key2, newBlock.Hash...)
		_ = txn.Set(newBlock.Hash, newBlock.Serialize())
		_ = txn.Set(key2, newBlock.Serialize())
		_ = txn.Set(key, newBlock.Hash)
		return nil
	})
	return nil
}

func (s basicService) PrintBlockchain(ctx context.Context, Hash string) (*BlockchainIter, error){
	var tip []byte
	key := append([]byte(Hash),[]byte("-last")...)


	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			log.Print("Initialising New Blockchain", err)
		} else {
			tip,_ = item.ValueCopy(nil)
		}
		return nil
	})
	if err != nil{
		return nil, err
	}
	bc := Blockchain{tip}
	log.Print("Failed here")
	blocks := &BlockchainIter{[][]byte{}}
	bci := bc.Iterator()
	for {
		block := bci.Next(s.db, Hash)
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






