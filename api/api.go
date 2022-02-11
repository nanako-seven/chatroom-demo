package api

type ServerChatMessageType int

const (
	Message ServerChatMessageType = iota
	Enter
	Leave
)

type ClientUserEnter struct {
	Username string
}

type ServerUserEnter struct {
	OK bool
}

type ClientChatMessage struct {
	Content string
}

type ServerChatMessage struct {
	Username string
	Type     ServerChatMessageType
	Content  string
}
