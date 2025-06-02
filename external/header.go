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

import "strconv"

type Header struct {
	BlockHeight uint64 `json:"block_height"`
	BlockTime   int64  `json:"block_time"`
	BlockHash   []byte `json:"hash"`
	ParentSlot  uint64 `json:"parent_slot"`
	ParentHash  []byte `json:"parent_hash"`
	Version     int8   `json:"version"`
}

func (h *Header) BlockHeightString() string {
	return strconv.FormatUint(h.BlockHeight, 10)
}
func (h *Header) BlockHeightBytes() []byte {
	return []byte(h.BlockHeightString())
}
