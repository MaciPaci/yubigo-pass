package utils

// Session holds active user session data, including identifiers and cryptographic material.
type Session struct {
	userID     string
	passphrase string
	salt       string
}

// NewEmptySession returns a new Session instance with all fields cleared,
// representing an unauthenticated state.
func NewEmptySession() Session {
	return Session{
		userID:     "",
		passphrase: "",
		salt:       "",
	}
}

// NewSession creates a new Session instance populated with the provided user data.
// This typically occurs after successful authentication.
func NewSession(userID string, passphrase string, salt string) Session {
	return Session{
		userID:     userID,
		passphrase: passphrase,
		salt:       salt,
	}
}

// Clear resets the session fields to their empty values, effectively logging the user out.
func (s *Session) Clear() {
	s.userID = ""
	s.passphrase = ""
	s.salt = ""
}

// GetUserID returns the unique identifier of the logged-in user.
// Returns an empty string if the session is not authenticated.
func (s Session) GetUserID() string {
	return s.userID
}

// GetPassphrase returns the raw passphrase stored in the session.
// Handle this value securely. Returns an empty string if the session is not authenticated.
func (s Session) GetPassphrase() string {
	return s.passphrase
}

// GetSalt returns the user-specific salt stored in the session.
// Returns an empty string if the session is not authenticated.
func (s Session) GetSalt() string {
	return s.salt
}

// IsAuthenticated checks if the session represents a logged-in user.
// It currently checks if the UserID field is non-empty.
func (s Session) IsAuthenticated() bool {
	return s.userID != ""
}

// GetSessionData returns all core session fields: userID, passphrase, and salt.
// Use with caution due to the sensitive nature of the passphrase.
func (s Session) GetSessionData() (string, string, string) {
	return s.userID, s.passphrase, s.salt
}
