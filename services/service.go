package services


import (
	"context"
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160" 
	"golang.org/x/crypto/sha3"
	"encoding/hex"
	"crypto/x509"
	"encoding/pem"
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

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}	



// func ValidateAddress(address string) bool {
// 	pubKeyHash := Base58Decode([]byte(address))
// 	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
// 	version := pubKeyHash[0]
// 	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
// 	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

// 	return bytes.Compare(actualChecksum, targetChecksum) == 0
// }


func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}


func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}

	return decode
}

func encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
    x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
    pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

    x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
    pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

    return string(pemEncoded), string(pemEncodedPub)
}

func decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
    block, _ := pem.Decode([]byte(pemEncoded))
    x509Encoded := block.Bytes
    privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

    blockPub, _ := pem.Decode([]byte(pemEncodedPub))
    x509EncodedPub := blockPub.Bytes
    genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
    publicKey := genericPublicKey.(*ecdsa.PublicKey)

    return privateKey, publicKey
}

func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
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

type Wallet struct {
	PrivateKey string `json:"privateKey"`
	PubKey string `json:"pubKey"`
	Address string `json:"address"`
}




