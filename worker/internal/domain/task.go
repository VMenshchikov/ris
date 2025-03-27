package domain

type Task struct {
	OrderId     uint
	TargetHash  [16]byte
	MaxLen      uint
	BlockNumber uint
	BlockSize   uint
}
