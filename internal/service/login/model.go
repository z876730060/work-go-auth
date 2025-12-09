package login

import "errors"

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l LoginReq) Validate() error {
	if l.Username == "" {
		return errors.New("username is required")
	}
	if l.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
