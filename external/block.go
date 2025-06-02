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
