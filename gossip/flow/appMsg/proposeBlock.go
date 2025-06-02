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
	_, err := f.domain.VerifyInsertBlockAndTransaction(self.Block)
	if err != nil {
		log.Fatalf("VerifyInsertBlockAndTransaction err: %s", err)
	}

}
