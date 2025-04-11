package dto

type Task struct {
	OrderId     uint64   `json:"orderId"`
	TargetHash  [16]byte `json:"targetHash"`
	BlockSize   uint     `json:"blockSize"`
	BlockNumber uint     `json:"blockNumber"`
	MaxLen      uint     `json:"maxLen"`
}
