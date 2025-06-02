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
