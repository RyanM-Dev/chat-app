package domain

type Message struct {
	ID       ID
	SenderID ID
	ChatID   ID
	Content  string
}
