package external

import "time"

type NodeStatus struct {
	SyncedTime    time.Time
	Synced        bool
	LastUpdated   time.Time
	IsConnected   bool
	StartTime     time.Time
	LastBlockTime time.Time
	LastBlockHash []byte
	LastHeight    uint64
}

func NewNodeStatus() *NodeStatus {
	return &NodeStatus{
		Synced:      false,
		IsConnected: false,
		StartTime:   time.Now(),
	}
}
