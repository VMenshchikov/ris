package model

import "time"

type Worker struct {
	Id         int64     `bson:"worker_id"`
	MaxTasks   uint      `bson:"max_tasks"`
	LastAction time.Time `bson:"last_action"`
}
