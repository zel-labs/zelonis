package appMsg

import (
	"encoding/json"
)

type Flow struct {
	Header  int    `json:"header"`
	Payload []byte `json:"payload"`
}

func NewFlow() *Flow {
	return &Flow{}
}

func (f *Flow) Encode() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Flow) Decode(b []byte) error {
	return json.Unmarshal(b, f)
}

func (f *Flow) FilterPayload() {
	switch f.Header {
	case SendInvBlockHash:
		payload := NewInvBlockHash()
		payload.Decode(f.Payload)
		
	}
}
