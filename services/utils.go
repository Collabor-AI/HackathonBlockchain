package services

import (
    "crypto/sha256"
    "golang.org/x/crypto/ripemd160" 
    "github.com/mr-tron/base58"
    "golang.org/x/crypto/sha3"
    "crypto/ecdsa"
    "log"
    "crypto/elliptic"
    "crypto/x509"
    "encoding/pem"
    "crypto/rand"
)
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
//  pubKeyHash := Base58Decode([]byte(address))
//  actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
//  version := pubKeyHash[0]
//  pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
//  targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

//  return bytes.Compare(actualChecksum, targetChecksum) == 0
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
    log.Print("this5")
    genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
    log.Print("this6")
    publicKey := genericPublicKey.(*ecdsa.PublicKey)
    log.Print("this7")

    return privateKey, publicKey
}

func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}