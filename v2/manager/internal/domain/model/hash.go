package model

import (
	"time"
)

type OrderStatusType string
type TaskStatusType string

const (
	IN_PROGRESS OrderStatusType = "IN_PROGRESS"
	ALL_TASKS   OrderStatusType = "ALL_TASKS"
	READY       OrderStatusType = "READY"
	ERROR       OrderStatusType = "ERROR"

	CREATED   TaskStatusType = "CREATED"
	SENDED    TaskStatusType = "SENDED"
	COMPLETED TaskStatusType = "COMPLETED"
)

type Sequence struct {
	Name  string `bson:"name"`
	Value uint64 `bson:"value"`
}

type OrderInfo struct {
	Id         uint64          `bson:"id,omitempty"`
	Status     OrderStatusType `bson:"status"`
	TargetHash [16]byte        `bson:"target_hash"`
	MaxLen     uint            `bson:"max_len"`
	Timeout    time.Time       `bson:"timeout"`
	BlockSize  int64           `bson:"block_size"`
	Results    []string        `bson:"results"`
	CreatedAt  time.Time       `bson:"created_at"`
}

type TaskInfo struct {
	OrderId     uint64         `bson:"order_id"`
	Status      TaskStatusType `bson:"status"`
	BlockNumber uint           `bson:"block_number"`
	UpdatedAt   time.Time      `bson:"updated_at"`
}
