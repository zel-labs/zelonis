package appMsg

import (
	"encoding/json"
	"log"
	"zelonis/external"
)

type ProposeBlock struct {
	*external.Block
}

func NewProposeBlock() *ProposeBlock {
	return &ProposeBlock{}
}

func (self *ProposeBlock) Decode(data []byte) error {
	return json.Unmarshal(data, self)
}

func (self *ProposeBlock) Encode() ([]byte, error) {
	return json.Marshal(self)
}

func (self *ProposeBlock) Process(f *flowControl) {
	f.domain.VerifyInsertBlockAndTransaction(self.Block)
	log.Println("Block Proposed by node")
}
