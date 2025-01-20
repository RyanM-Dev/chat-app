package domain

import "time"

const (
	Admin  = "Admin"
	Owner  = "Owner"
	Normal = "Normal"
)

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
