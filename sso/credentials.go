package sso

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Credentials) Validate() error {
	if c.Username == "" || c.Password == "" {
		return EmptyCredsError
	}
	return nil
}
