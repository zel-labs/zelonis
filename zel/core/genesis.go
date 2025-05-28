package core

type Genesis struct {
	BlockHeight  []byte         `json:"block_height"`
	BlockTime    []byte         `json:"block_time"`
	BlockHash    []byte         `json:"hash"`
	ParentSlot   []byte         `json:"parent_slot"`
	ParentHash   []byte         `json:"parent_hash"`
	Transactions []*Transaction `json:"transactions"`
	Version      []byte         `json:"version"`
}

type Transaction struct {
	Inpoint   *Inpoint    `json:"inpoint"`
	Outpoints []*Outpoint `json:"outpoints"`
	Balance   []byte      `json:"balance"`
	TxHash    []byte      `json:"transaction_hash"`
	Signature []byte      `json:"signature"`
}

type Inpoint struct {
	PubKey        []byte `json:"pub_key"`
	Value         []byte `json:"value"`
	PrevBlockHash []byte `json:"prev_block_hash"`
}

type Outpoint struct {
	PubKey []byte `json:"pub_key"`
	Value  []byte `json:"value"`
}
