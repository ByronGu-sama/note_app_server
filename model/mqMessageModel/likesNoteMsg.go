package mqMessageModel

import "time"

type LikeNotes struct {
	Action    string
	Nid       string
	Uid       uint
	Timestamp time.Time
}
