package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/facundocarballo/basics-blockchain-go/handlers"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLenght = 4
	version        = byte(0x00) // Respect the length of the checksum, in this case is 4 so our version will be 4 bytes.
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w *Wallet) Address() []byte {
	publicHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, publicHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)

	fmt.Printf("Public Key:   0x%x\n", w.PublicKey)
	fmt.Printf("Public Hash:  0x%x\n", publicHash)
	fmt.Printf("Address:      %s\n", Base58Encode(fullHash))

	return Base58Encode(fullHash)
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	handlers.HandleErrors(err)

	publicKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, publicKey
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	return &Wallet{PrivateKey: private, PublicKey: public}
}

func PublicKeyHash(pubKey []byte) []byte {
	publicKeyHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(publicKeyHash[:])
	handlers.HandleErrors(err)

	return hasher.Sum(nil)
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLenght] // This returns the first 4 bytes of the second hash.
}
