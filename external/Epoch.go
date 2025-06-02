package external

type Epoch struct {
	EpochNumber uint64 `json:"epoch_number"`
	EpochStart  uint64 `json:"epoch_start"`
	EpochEnd    uint64 `json:"epoch_end"`
	EpochReward uint64 `json:"epoch_reward"`
}
