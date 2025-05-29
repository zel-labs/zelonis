package external

import (
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

const DomainHashSize = 32

type DomainTransactionID DomainHash

type DomainHash struct {
	hashArray [DomainHashSize]byte
}

func NewDomainHashFromByteArray(hashBytes *[DomainHashSize]byte) *DomainHash {
	return &DomainHash{
		hashArray: *hashBytes,
	}
}

func NewDomainHashFromByteSlice(hashBytes []byte) (*DomainHash, error) {
	if len(hashBytes) != DomainHashSize {
		return nil, errors.Errorf("invalid hash size. Want: %d, got: %d",
			DomainHashSize, len(hashBytes))
	}
	domainHash := DomainHash{
		hashArray: [DomainHashSize]byte{},
	}
	copy(domainHash.hashArray[:], hashBytes)
	return &domainHash, nil
}

func NewDomainHashFromString(hashString string) (*DomainHash, error) {
	expectedLength := DomainHashSize * 2
	// Return error if hash string is too long.
	if len(hashString) != expectedLength {
		return nil, errors.Errorf("hash string length is %d, while it should be be %d",
			len(hashString), expectedLength)
	}

	hashBytes, err := hex.DecodeString(hashString)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return NewDomainHashFromByteSlice(hashBytes)
}

func EncryptAndHash(inter interface{}) (*DomainHash, error) {
	encrypted, err := json.Marshal(inter)
	if err != nil {
		return nil, err
	}

	Txhash := blake2b.Sum256(encrypted)
	return NewDomainHashFromByteArray(&Txhash), nil
}

func (hash DomainHash) String() string {
	return hex.EncodeToString(hash.hashArray[:])
}

func (hash *DomainHash) ByteArray() *[DomainHashSize]byte {
	arrayClone := hash.hashArray
	return &arrayClone
}

func (hash *DomainHash) ByteSlice() []byte {
	return hash.ByteArray()[:]
}

func NewDomainTransactionIDFromByteSlice(transactionIDBytes []byte) (*DomainTransactionID, error) {
	hash, err := NewDomainHashFromByteSlice(transactionIDBytes)
	if err != nil {
		return nil, err
	}
	return (*DomainTransactionID)(hash), nil
}
