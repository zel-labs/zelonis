/*
Copyright (C) 2025 Zelonis Contributors

This file is part of Zelonis.

Zelonis is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Zelonis is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Zelonis. If not, see <https://www.gnu.org/licenses/>.
*/
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
