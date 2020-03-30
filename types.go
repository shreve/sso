package sso

type User struct {
	Uid string
	Username string
	Password string `json:"-"`
}

type Token struct {
	Token string `json:"token"`
	User *User `json:"-"`
}

type Credentials struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
}

func (c *Credentials) Validate() error {
	if c.Username == "" || c.Password == "" {
		return EmptyCredsError
	}
	return nil
}
