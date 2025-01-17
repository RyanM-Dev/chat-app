package domain

import "time"

type Chat struct {
	ID          ID
	Name        string
	Owner       ID
	Admins      []ID
	Members     []ID
	CreatedTime *time.Time
	DeletedTime *time.Time
	ChatType    ChatType
}
