package domain

type Session struct {
	SessionID     ID
	UserID        ID
	ChatIDAndName map[string]string // Map to store chat membership with roles (user, admin, owner)
	ChatNameList  []string
}
