package model

// Password is the model of the password
type Password struct {
	UserID   string `db:"user_id"`
	Title    string `db:"title"`
	Username string `db:"username"`
	Password string `db:"password"`
	Url      string `db:"url"`
	Nonce    []byte `db:"nonce"`
}

// NewPassword returns new Password instance
func NewPassword(userID, title, username, password, url string, nonce []byte) Password {
	return Password{
		UserID:   userID,
		Title:    title,
		Username: username,
		Password: password,
		Url:      url,
		Nonce:    nonce,
	}
}
