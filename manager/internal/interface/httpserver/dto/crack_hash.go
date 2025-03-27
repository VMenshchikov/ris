package dto

type CrackHashRequest struct {
	Hash      string `json:"hash"`
	MaxLength uint   `json:"maxLength"`
	Timeout   uint   `json:"timeout"`
	Weight    uint   `json:"weight"`
	BlockSize uint   `json:"blockSize"`
}

func (c *CrackHashRequest) SetDefaults() bool {
	if c.Hash == "" || c.MaxLength == 0 {
		return false
	}

	if c.Timeout == 0 {
		c.Timeout = 60
	}

	if c.Weight == 0 {
		c.Weight = 1
	}

	if c.BlockSize == 0 {
		c.BlockSize = 1_000_000
	}
	return true
}

type CrackHashResponse struct {
	Id uint64 `json:"requestId"`
}

type GetResultResponse struct {
	Status   string   `json:"status"`
	Results  []string `json:"results"`
	Progress float64  `json:"progress"`
}

type WorkerResult struct {
	WorkerId   uint     `json:"workerId"` //todo to string
	OrderId    uint64   `json:"orderId"`
	TaskNumber uint     `json:"number"`
	Results    []string `json:"results"`
}
