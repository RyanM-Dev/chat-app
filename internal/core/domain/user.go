package domain

import "time"

type Gender int

const (
	Male Gender = iota + 1
	Female
	NonBinary
)

type ChatType int

const (
	Private ChatType = iota + 1
	Group
)

type ID string

type User struct {
	ID          ID
	Username    string
	FirstName   string
	LastName    string
	Password    string
	Gender      Gender
	Email       string
	Contacts    []ID
	DateOfBirth *time.Time
	CreatedTime *time.Time
	DeletedTime *time.Time
}
