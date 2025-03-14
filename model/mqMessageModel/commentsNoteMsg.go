package mqMessageModel

import (
	"encoding/json"
	"time"
)

type DelNoteComment struct {
	Action    string
	Cid       string
	Uid       uint
	Timestamp time.Time
}

func (msg *DelNoteComment) Decode(object []byte) error {
	return json.Unmarshal(object, msg)
}

func (msg *DelNoteComment) Encode() ([]byte, error) {
	return json.Marshal(msg)
}
