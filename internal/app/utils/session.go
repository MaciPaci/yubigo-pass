package utils

type Session struct {
	userID     string
	passphrase string
	salt       string
}

func NewEmptySession() Session {
	return Session{
		userID:     "",
		passphrase: "",
		salt:       "",
	}
}

func NewSession(userID string, passphrase string, salt string) Session {
	return Session{
		userID:     userID,
		passphrase: passphrase,
		salt:       salt,
	}
}

func (s *Session) Clear() {
	s.userID = ""
	s.passphrase = ""
	s.salt = ""
}

func (s Session) GetUserID() string {
	return s.userID
}

func (s Session) GetPassphrase() string {
	return s.passphrase
}

func (s Session) GetSalt() string {
	return s.salt
}

func (s Session) GetSessionData() (string, string, string) {
	return s.userID, s.passphrase, s.salt
}
