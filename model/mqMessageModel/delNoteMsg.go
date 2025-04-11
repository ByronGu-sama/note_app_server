package mqMessageModel

import (
	"encoding/json"
	"time"
)

type DelNote struct {
	Action    int
	Uid       int64
	Nid       string
	Timestamp time.Time
}

func (that *DelNote) Encode() ([]byte, error) {
	return json.Marshal(that)
}

func (that *DelNote) Decode(bts []byte) error {
	return json.Unmarshal(bts, that)
}
