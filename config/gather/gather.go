package gather

import (
	"encoding/json"
	"time"
)

const timeFormat = `2006-01-02T15:04:05.999999999Z`

//采集数据的结构
type GatherData struct {
	Timestamp time.Time              `json:"timestamp"`
	Message   string                 `json:"message"`
	Tags      []string               `json:"tags,omitempty"`
	Extra     map[string]interface{} `json:"-"`
}

func appendtotags(tags []string, tag string) []string {
	for _, t := range tags {
		if t == tag {
			return tags
		}
	}
	return append(tags, tag)
}

//添加标签
func (self *GatherData) AddTag(tags ...string) {
	for _, tag := range tags {
		self.Tags = append(self.Tags, tag)
	}
}

func (self *GatherData) conver2Map() map[string]interface{} {
	gatherMap := map[string]interface{}{
		"@timestamp": self.Timestamp.UTC().Format(timeFormat),
	}
	if self.Message != "" {
		gatherMap["message"] = self.Message
	}

	if len(self.Tags) > 0 {
		gatherMap["tags"] = self.Tags
	}

	for key, value := range self.Extra {
		gatherMap[key] = value
	}
	return gatherMap
}

func (self *GatherData) MarshalJSON() (data []byte, err error) {
	gm := self.conver2Map()
	return json.Marshal(gm)
}

func (self *GatherData) MarshalIndent() (data []byte, err error) {
	gm := self.conver2Map()
	return json.MarshalIndent(gm, "", "\t")
}
