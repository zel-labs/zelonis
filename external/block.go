package external

import (
	"encoding/json"
	"time"
)

type Block struct {
	Header       *Header        `json:"h"`  //header
	Transactions []*Transaction `json:"t"`  //transactions
	Validator    *ValidatorInfo `json:"v"`  //validator
	Signature    []byte         `json:"s"`  //signature
	RecivedAt    time.Time      `json:"r"`  //recivedTime
	RecivedFrom  []byte         `json:"rf"` //recievedfrom
}

func (b *Block) Serialize() []byte {

	info, _ := json.Marshal(b)
	return info
}
func (b *Block) Deserialize(data []byte) error {

	err := json.Unmarshal(data, b)
	return err
}
