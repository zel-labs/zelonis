package appMsg

import (
	"encoding/json"
	"zelonis/external"
)

type BlockInfo struct {
	*external.Block
}

func NewBlockInfo() *BlockInfo {
	return &BlockInfo{}
}

func (self *BlockInfo) Decode(data []byte) error {
	return json.Unmarshal(data, self)
}

func (self *BlockInfo) Encode() ([]byte, error) {
	return json.Marshal(self)
}
