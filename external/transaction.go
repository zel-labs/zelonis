package external

import "encoding/json"

type Transaction struct {
	Inpoint   *Inpoint    `json:"inpoint"`
	Outpoints []*Outpoint `json:"outpoints"`
	TxHash    []byte      `json:"transaction_hash"`
	Signature []byte      `json:"signature"`
	Fee       []byte      `json:"fee"`
	Status    int8        `json:"status"`
	TxType    int8        `json:"tx_type"`
}

type Inpoint struct {
	PubKey        []byte `json:"pub_key"`
	Value         []byte `json:"value"`
	PrevBlockHash []byte `json:"prev_block_hash"`
}
type Outpoint struct {
	PubKey []byte `json:"pub_key"`
	Value  []byte `json:"value"`
	TxType int8   `json:"tx_type"`
}

func (tx *Transaction) TxSerialize() []byte {
	txBytes, _ := json.Marshal(tx)
	return txBytes
}

func (tx *Transaction) DbTxToDomainTX(txByte []byte) {
	json.Unmarshal(txByte, tx)

}
