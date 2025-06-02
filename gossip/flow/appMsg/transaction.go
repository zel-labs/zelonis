package appMsg

import (
	"encoding/json"
	"zelonis/external"
)

type InviTransaction struct {
	*external.Transaction
}

func NewInviTransaction() *InviTransaction {
	return &InviTransaction{}
}

func (self *InviTransaction) Decode(data []byte) error {
	return json.Unmarshal(data, self)
}

func (self *InviTransaction) Encode() ([]byte, error) {
	return json.Marshal(self)
}

func (self *InviTransaction) Process(f *flowControl) {
	status := f.domain.TxManager().Mempool().AddTxToMempool(self.Transaction)
	if status {
		f.encodeAndSend(self.Transaction, SendInviTransaction)
	}
}
