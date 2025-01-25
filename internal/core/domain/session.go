package domain

// ChatMapping defines a mapping between chat names and their IDs.
type ChatMapping map[string]ID

// Session represents a user session and associated chat information.
type Session struct {
	SessionID    ID // Unique identifier of a session.
	UserID       ID // Identifier of the user owning the session.
	ChatMappings ChatMapping
	ChatNames    []string // List of chat names associated with the session.
}
