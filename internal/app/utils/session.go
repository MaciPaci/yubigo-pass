package utils

// Session is a model for holding session data
type Session struct {
	userID     string
	passphrase string
	salt       string
}

// NewEmptySession returns new empty Session instance
func NewEmptySession() Session {
	return Session{
		userID:     "",
		passphrase: "",
		salt:       "",
	}
}

// NewSession returns new Session instance
func NewSession(userID string, passphrase string, salt string) Session {
	return Session{
		userID:     userID,
		passphrase: passphrase,
		salt:       salt,
	}
}

// Clear clears the current session
func (s *Session) Clear() {
	s.userID = ""
	s.passphrase = ""
	s.salt = ""
}

// GetUserID returns user ID from the session
func (s Session) GetUserID() string {
	return s.userID
}

// GetPassphrase returns passphrase from the session
func (s Session) GetPassphrase() string {
	return s.passphrase
}

// GetSalt returns salt from the session
func (s Session) GetSalt() string {
	return s.salt
}

// GetSessionData returns session data
func (s Session) GetSessionData() (string, string, string) {
	return s.userID, s.passphrase, s.salt
}
