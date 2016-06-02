package gather

import (
	"time"
)

//采集数据的结构
type GatherData struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}
