package dto

type ResultMessage struct {
	WorkerId   uint     `json:"workerId"`
	OrderId    uint64   `json:"orderId"`
	TaskNumber uint     `json:"number"`
	Results    []string `json:"results"`
}

type AliveMessage struct {
	WorkerId uint `json:"workerId"`
	MaxTasks uint `json:"maxTasks"`
}
