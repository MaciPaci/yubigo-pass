package model

// Password is the model of the password
type Password struct {
	ID       string `db:"id"`
	UserID   string `db:"user_id"`
	Title    string `db:"title"`
	Username string `db:"username"`
	Password string `db:"password"`
	Url      string `db:"url"`
	Nonce    []byte `db:"nonce"`
}

// NewPassword returns new Password instance
func NewPassword(id, userID, title, username, password, url string, nonce []byte) Password {
	return Password{
		ID:       id,
		UserID:   userID,
		Title:    title,
		Username: username,
		Password: password,
		Url:      url,
		Nonce:    nonce,
	}
}
