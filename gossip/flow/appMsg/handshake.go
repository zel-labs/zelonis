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
	"github.com/pkg/errors"
)

const (
	handShakeVersion = "0.0.01"
)

type handshake struct {
	Version string `json:"version"`
	Action  int    `json:"action"`
}

func NewHandshakeWithInfo() *handshake {
	return &handshake{
		Version: handShakeVersion,
		Action:  HandshakeMsg,
	}
}
func NewHandshake() *handshake {
	return &handshake{}
}

func (h *handshake) Decode(data []byte) error {
	return json.Unmarshal(data, h)
}
func (h *handshake) Encode() ([]byte, error) {
	return json.Marshal(h)
}

func (h *handshake) Verify(msg *handshake) error {
	if h.Version != msg.Version {

		return errors.New("handshake version mismatch")
	} else if h.Action != msg.Action || h.Action != HandshakeMsg {
		return errors.New("handshake action mismatch")
	}
	return nil
}
