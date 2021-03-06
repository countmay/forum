package session

import (
	uuid "github.com/satori/go.uuid"
)

// SessionData ...
type SessionData struct {
	Username     string
	IsAuthorized bool
	// Expiration time.time
}

// Session ...
type Session struct {
	data map[string]SessionData
}

// NewSession ...
func NewSession() *Session {
	s := new(Session)

	s.data = make(map[string]SessionData)

	return s
}

// Init ...
func (s *Session) Init(username string) (string, SessionData) {
	sessionID, _ := uuid.NewV4()

	data := SessionData{
		Username:     username,
		IsAuthorized: true,
	}
	s.data[sessionID.String()] = data

	return sessionID.String(), data
}

func (s *Session) Delete(token string) {
	delete(s.data, token)
}

func (s *Session) Data(token string) SessionData {
	return s.data[token]
}

func (s *Session) Authed(token string) bool {
	if _, ok := s.data[token]; !ok {
		return false
	}

	return true
}

func (s *Session) CheckUsersSession(username string) (bool, string) {
	for key, value := range s.data {
		if value.Username == username {
			return true, key
		}
	}
	return false, ""
}
