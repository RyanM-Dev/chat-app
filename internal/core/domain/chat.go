package domain

import "time"

type Chat struct {
	ID          ID
	Name        string
	Members     []ID
	CreatedTime time.Time
	DeletedTime time.Time
	ChatType    ChatType
}
