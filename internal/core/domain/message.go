package domain

import "time"

type Message struct {
	ID          ID
	SenderID    ID
	ChatID      ID
	CreatedTime *time.Time
	Content     string
}
