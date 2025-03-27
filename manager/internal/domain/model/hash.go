package model

import (
	"sync"
	"time"
)

type StatusType string

const (
	IN_PROGRESS StatusType = "IN_PROGRESS"
	ALL_TASKS   StatusType = "ALL_TASKS"
	READY       StatusType = "READY"
	ERROR       StatusType = "ERROR"
)

type OrderInfo struct {
	Id         uint64
	Status     StatusType
	TargetHash [16]byte
	MaxLen     uint
	Timeout    time.Time
	Weight     uint
	BlockSize  int64
	Results    []string

	TasksUncompleted uint
	Lock             sync.Mutex
}
