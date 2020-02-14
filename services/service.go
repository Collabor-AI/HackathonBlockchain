package services


import (
	"context"
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	firebase "firebase.google.com/go"
	firebaseDB "firebase.google.com/go/db"
	//github.com/google/uuid
	// "github.com/go-kit/kit/log"
	"log"
	// "strconv"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

// Service  methods
type Service interface {
	NewBlockchain(ctx context.Context, startingData InitData) (Blockchain, error)
	AddBlock(ctx context.Context, bd BlockData, poml float64, PubKey string ,PrivKey string, Hash string) (error)
	PrintBlockchain(ctx context.Context, Hash string) (*BlockchainIter, error)
	GenerateAddress(ctx context.Context)  (*Wallet, error)
	PrintLeaderBoard(ctx context.Context, Hash string)  (map[string]gbms, error)
}

//Our Service contains badger, firebase instance, and firebase DB
func New(db *badger.DB, app *firebase.App, client *firebaseDB.Client) Service {
	var svc Service
	{
		svc = NewBasicService(db, app, client)
	}
	return svc
}	

func NewBasicService(db *badger.DB, app *firebase.App, client *firebaseDB.Client) Service {
	return basicService{db, app, client}
}

type basicService struct{
	//For now our service only needs the blockchain - which is in persistence, badger DB
	db *badger.DB
	app *firebase.App
	client *firebaseDB.Client
}

func (s basicService) NewBlockchain(ctx context.Context, startingData InitData) (Blockchain, error) {
	//Check if there is an existing blockchain
	var tip []byte
	var base64tip string
	key := append([]byte(startingData.Dataset.Name),[]byte("-last")...)
	
		// log.Print("Key: %+v",key)X
		// err := s.db.View(func(txn *badger.Txn) error {
		// 	item, err := txn.Get(key)
		// 	if err != nil {
		// 		log.Print("Initialising New Blockchain", err)
		// 	} else {
		// 		tip,_ = item.ValueCopy(nil)
		// 	}
		// 	return nil
		// })
	var err error
	if err = s.client.NewRef(string(key)).Get(ctx, &base64tip); err != nil {
	  log.Print(err)
	  fmt.Print("Not foundn inn Firebase\n")
	} else {
		fmt.Printf("Tip of the blockchain is %v",base64tip)
		tip, err = base64.StdEncoding.DecodeString(base64tip)
	}
	

 
	var hashList []string
	if len(tip) == 0 {
		genesis := NewGenesisBlock(startingData)
		blockData := genesis.Serialize()

		/* Badger Write
		// _ = s.db.Update(func(txn *badger.Txn) error {		

		// 	// Set "{name}-last" to genesisHash
		// 	_ = txn.Set(key, genesis.Hash)
			

		// 	// Set "Hash" to BlockData
		// 	blockData := genesis.Serialize()
		// 	_ = txn.Set(genesis.Hash, blockData)
			

		// 	//Set "hashList" to hashList + [genesisHash]
		// 	item2, err2 := txn.Get([]byte("hashList"))			
		// 	if err2 != nil {
		// 		log.Print("Failed here1")
		// 		hashList = append(hashList,[][]byte{genesis.Hash}...)
		// 	} else {
		// 		data,_ := item2.ValueCopy(nil)
		// 		_ = json.Unmarshal(data,&hashList)
		// 		hashList = append(hashList,[][]byte{genesis.Hash}...)
		// 	}
			
		// 	hashListBytes,_ := json.Marshal(hashList)
		// 	_ = txn.Set([]byte("hashList"), hashListBytes)
		// 	tip = genesis.Hash
		// 	log.Print("genesis is %+v",genesis)
		// 	return nil
		// })
		*/

		// Set "{name}-last" to genesisHash
		if err := s.client.NewRef(string(key)).Set(ctx, base64.StdEncoding.EncodeToString(genesis.Hash)); err != nil {
		  log.Print(err)
		  fmt.Print("1: Failed to set name-last : genesisHash to firebase\n")
		}

		//Set genesisHash : blockData		
		if err = s.client.NewRef(base64.StdEncoding.EncodeToString(genesis.Hash)).Set(ctx, blockData); err != nil {
			  log.Print(err)
			  fmt.Print("2: Failed to set genesisHash : blockData to firebase\n")
		}

		//get hashlist
		// var hashList [][]byte 
		if err = s.client.NewRef("hashList").Get(ctx, &hashList); err != nil {
		  log.Print(err)
		  fmt.Print("HashList not found in Firebase\n")
		}

		// _ = json.Unmarshal(hashListBytesFromDB,&hashList)
		hashList = append(hashList,[]string{startingData.Dataset.Name}...)
		// hashListBytes,_ := json.Marshal(hashList)
		//Set "hashList" to hashList + [genesisHash]
		if err := s.client.NewRef("hashList/").Set(ctx, hashList); err != nil {
			  log.Print(err)
			  fmt.Print("Failed to write to firebase2\n")
		}
		fmt.Printf("HashList: %+v",hashList)
		tip = genesis.Hash
	}



	bc := Blockchain{tip}
	if err != nil {
		log.Print("failed to initiate blockchain")
	}
	log.Print("BC %+v",bc)
	
	//return the tip of the blockchain
	return bc, nil
}

type gbms struct {
	Hash string `json:"hash"`
	Score float64 `json:"score"`
	Email string `json:"email"`
	Timestamp int64
}



func (s basicService) AddBlock(ctx context.Context, bd BlockData, poml float64, PubKey string , PrivKey string, Hash string) error{
	var lastHash string
	key := append([]byte(Hash),[]byte("-last")...)


	// // retrieve end of blockchain from badger
	// _ = s.db.View(func(txn *badger.Txn) error {
	// 	item, _ := txn.Get(key)
	// 	lastHash, _ = item.ValueCopy(nil)
	// 	return nil
	// })


	//get name-lastHash from firebase
	if err := s.client.NewRef(string(key)).Get(ctx, &lastHash); err != nil {
	  log.Print(err)
	  fmt.Print("1: Failed to get lastHash from firebase\n")
	}
		
	//Validate address
	lastHashBytes, _ := base64.StdEncoding.DecodeString(lastHash)
	fmt.Print("lastHash (string): %v",string(lastHash))
	fmt.Print("lastHash (string): %v",string(lastHashBytes))
	log.Printf("PrivKey %v PubKey %v",PrivKey ,PubKey)
	_,pub := decode(PrivKey,PubKey)
    pubKeyBytes := elliptic.Marshal(elliptic.P384(), pub.X, pub.Y)
    address := "0x" + hex.EncodeToString(Keccak256(pubKeyBytes[1:])[12:]) 
	bd.Address = address

	//Get Global Model Set
	var gbmsData map[string]gbms;
	if err := s.client.NewRef("gbms/"+Hash+"/").Get(ctx, &gbmsData); err != nil {
	  log.Print(err)
	  fmt.Print("1: Failed to get GBMS from firebase\n")
	}
	write := true

	// only add if its better than poml or address not in 

	if val, ok := gbmsData[address]; ok {
		if val.Score >= poml {
			write = false
		}
	}
	if write {
		//create a new block
		bdBytes,_ := json.Marshal(bd)
		newBlock := NewBlock(bdBytes, poml, lastHashBytes)

		//write to DB
		// _ = s.db.Update(func(txn *badger.Txn) error {
		// 	_ = txn.Set(newBlock.Hash, newBlock.Serialize())
		// 	_ = txn.Set(key, newBlock.Hash)
		// 	return nil
		// })


		//Set Hash : blockData		
		if err := s.client.NewRef(base64.StdEncoding.EncodeToString(newBlock.Hash)).Set(ctx, newBlock); err != nil {
			  fmt.Print("2: Failed to set Hash : blockData to firebase\n")
		}
		//Set lastHash : Hash
		if err := s.client.NewRef(string(key)).Set(ctx,base64.StdEncoding.EncodeToString(newBlock.Hash)); err != nil {
			  fmt.Print("2: Failed to set Hash : blockData to firebase\n")
		}
		if len(gbmsData) == 0 {
			gbmsData = make(map[string]gbms)
		}  
			
		gbmsData[address] = gbms{base64.StdEncoding.EncodeToString(newBlock.Hash), poml, bd.Email, newBlock.Timestamp}
		
			
		
		
		
		if err := s.client.NewRef("gbms/"+Hash+"/").Set(ctx, gbmsData); err != nil {
		  log.Print(err)
		  fmt.Print("1: Failed to get GBMS from firebase\n")
		}
		gbmsDataBytes,_ := json.Marshal(gbmsData)
		if err := s.client.NewRef("gbms/"+Hash+"/").Set(ctx, gbmsDataBytes); err != nil {
		  log.Print(err)
		  fmt.Print("1: Failed to get GBMS from firebase\n")
		}
	}

	return nil
}

func (s basicService) PrintLeaderBoard(ctx context.Context, Hash string) (map[string]gbms, error){
	
	var gbmsDataBytes []byte
	var gbmsData map[string]gbms
	var err error
	if err = s.client.NewRef("gbms/"+Hash+"/").Get(ctx, &gbmsDataBytes); err != nil {
	  log.Print(err)
	  fmt.Print("Not foundn inn Firebase\n")
	}
	fmt.Printf("GBMS : %+v", gbmsDataBytes)
	_ = json.Unmarshal(gbmsDataBytes, &gbmsData)

	fmt.Printf("GBMS : %+v", gbmsData)
	// return nil, nil
	return gbmsData, err
}


func (s basicService) PrintBlockchain(ctx context.Context, Hash string) (*BlockchainIter, error){
	var tip []byte
	key := append([]byte(Hash),[]byte("-last")...)

	/*
	//Get Key from BC
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
	*/

	var base64tip string
	if err := s.client.NewRef(string(key)).Get(ctx, &base64tip); err != nil {
	  log.Print(err)
	  fmt.Print("Not foundn inn Firebase\n")
	}
	tip, _ = base64.StdEncoding.DecodeString(base64tip)


	fmt.Printf("Blockchain TIP: %v",base64tip)
	bc := Blockchain{tip}	
	blocks := &BlockchainIter{[][]byte{}}
	bci := bc.Iterator()
	for {
		block := bci.Next(s.db, s.client, ctx, Hash)
		blockData,_ := json.Marshal(block)
		blocks.Blocks = append(blocks.Blocks,blockData)
		if len(block.PrevBlockHash) == 0 {
			// fmt.Printf("%+v",block)
			// fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
			// fmt.Printf("Data: %+v\n", block.Data)
			// fmt.Printf("Hash: %x\n", block.Hash)
			break
		}

		

		// fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		// fmt.Printf("Data: %s\n", block.Data)
		// fmt.Printf("Hash: %x\n", block.Hash)
		// fmt.Println()

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






