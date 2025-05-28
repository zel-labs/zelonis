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
