package domain

type Session struct {
	SessionID     ID
	UserID        ID
	ChatNameAndID map[string]string // map[chatName]chatID
	ChatNameList  []string
}
