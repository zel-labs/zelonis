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

import "time"

type NodeStatus struct {
	SyncedTime    time.Time `json:"synced_time"`
	Synced        bool      `json:"synced,omitempty"`
	LastUpdated   time.Time `json:"last_updated"`
	IsConnected   bool      `json:"is_connected,omitempty"`
	StartTime     time.Time `json:"start_time"`
	LastBlockTime time.Time `json:"last_block_time"`
	LastBlockHash []byte    `json:"last_block_hash,omitempty"`
	LastHeight    uint64    `json:"last_height,omitempty"`
}

func NewNodeStatus() *NodeStatus {
	return &NodeStatus{
		Synced:      false,
		IsConnected: false,
		StartTime:   time.Now(),
	}
}
