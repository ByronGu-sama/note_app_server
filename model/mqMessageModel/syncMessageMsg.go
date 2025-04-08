package mqMessageModel

import (
	"encoding/json"
	"note_app_server/model/msgModel"
	"time"
)

type SyncMessageMsg struct {
	Action    string
	FirstKey  int64
	SecondKey int64
	Message   *msgModel.Message
	Timestamp time.Time
}

func (that *SyncMessageMsg) Encode() ([]byte, error) {
	return json.Marshal(that)
}

func (that *SyncMessageMsg) Decode(bts []byte) error {
	return json.Unmarshal(bts, that)
}
