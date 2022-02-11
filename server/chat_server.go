package main

import (
	"errors"
	"sync"
)

type ChatMessageType int

const (
	ChatMessageTypeMessage ChatMessageType = iota
	ChatMessageTypeEnter
	ChatMessageTypeLeave
)

type ChatServerMessage struct {
	Username string
	Type     ChatMessageType
	Content  string
}

type ChatServer struct {
	init  bool
	users map[string]chan ChatServerMessage
	mutex sync.RWMutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		users: make(map[string]chan ChatServerMessage),
		init:  true,
	}
}

func (s *ChatServer) UserEnter(username string) (chan ChatServerMessage, error) {
	if !s.init {
		panic("ChatServer not initialized")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.users[username]; ok {
		return nil, errors.New("username already exists")
	}
	msg := ChatServerMessage{
		Username: username,
		Type:     ChatMessageTypeEnter,
	}
	for _, v := range s.users {
		v <- msg
	}
	ch := make(chan ChatServerMessage)
	s.users[username] = ch
	return ch, nil
}

func (s *ChatServer) Broadcast(username string, content string) {
	if !s.init {
		panic("ChatServer not initialized")
	}
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	msg := ChatServerMessage{
		Username: username,
		Type:     ChatMessageTypeMessage,
		Content:  content,
	}
	for k, v := range s.users {
		if k != username {
			v <- msg
		}
	}
}

func (s *ChatServer) UserLeave(username string) {
	if !s.init {
		panic("ChatServer not initialized")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.users, username)
	msg := ChatServerMessage{
		Username: username,
		Type:     ChatMessageTypeLeave,
	}
	for _, v := range s.users {
		v <- msg
	}
}
