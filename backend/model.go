package backend

import "time"

type UIModel struct {
	Count  int
	Width  int
	Height int
	Time   time.Time
	Term   string
}
