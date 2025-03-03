package mqMessageModel

import (
	"encoding/json"
	"note_app_server/model/noteModel"
	"time"
)

type SyncNoteMsg struct {
	Action    string
	Note      *noteModel.ESNote
	Timestamp time.Time
}

func (that *SyncNoteMsg) Encode() ([]byte, error) {
	return json.Marshal(that)
}

func (that *SyncNoteMsg) Decode(bts []byte) error {
	return json.Unmarshal(bts, that)
}
