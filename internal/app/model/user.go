package model

// User is the model of the user
type User struct {
	Uuid     string
	Username string
	Password string
	Salt     string
}

// NewUser returns new User instance
func NewUser(uuid, username, password, salt string) User {
	return User{
		Uuid:     uuid,
		Username: username,
		Password: password,
		Salt:     salt,
	}
}
