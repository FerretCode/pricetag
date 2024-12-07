package session

import (
	"errors"
	"sync"
	"time"
)

type Session struct {
	UserID    int
	ExpiresAt time.Time
}

type SessionManager struct {
	sessions   sync.Map
	expiration time.Duration
}

func NewSessionManager(expiration time.Duration) *SessionManager {
	manager := &SessionManager{
		expiration: expiration,
	}

	go manager.cleanupExpiredSessions()
	return manager
}

func (sm *SessionManager) CreateSession(userID int) string {
	sessionID := generateSessionID()
	sm.sessions.Store(sessionID, Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(sm.expiration),
	})
	return sessionID
}

func (sm *SessionManager) GetSession(sessionID string) (Session, error) {
	value, ok := sm.sessions.Load(sessionID)
	if !ok {
		return Session{}, errors.New("session not found please log in")
	}

	session := value.(Session)
	if time.Now().After(session.ExpiresAt) {
		sm.sessions.Delete(sessionID)
		return Session{}, errors.New("session expired")
	}

	return session, nil
}

func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.sessions.Delete(sessionID)
}

func (sm *SessionManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sm.sessions.Range(func(key, value any) bool {
			session := value.(Session)
			if time.Now().After(session.ExpiresAt) {
				sm.sessions.Delete(key)
			}
			return true
		})
	}
}

func generateSessionID() string {
	return time.Now().Format("20060102150405.000000")
}
