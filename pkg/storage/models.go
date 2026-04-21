package storage

import "time"

type Available_Keys struct {
	Id        int64
	ShortKey  string
	CreatedAt time.Time
}

type Links struct {
	Id           int64
	ActualURL    string
	ShortKey     string
	CreatedAt    time.Time
	ExpiredAt    time.Time
	ClickCounter int
}
