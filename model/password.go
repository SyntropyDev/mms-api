package model

import (
	"fmt"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/dchest/uniuri"
)

type Password struct {
	text string
	hash string
}

func NewPassword(s string) (*Password, error) {
	if !isValidPassword(s) {
		err := fmt.Errorf("auth: invalid password - %v password must be between 7 and 32 characters", s)
		return nil, err
	}
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Password{s, string(b)}, nil
}

func NewAutoPassword() *Password {
	s := uniuri.NewLen(15)
	b, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return &Password{s, string(b)}
}

func (p *Password) Hash() string {
	return p.hash
}

func (p *Password) String() string {
	return p.text
}

func isValidPassword(p string) bool {
	hasMin := len(p) >= 7
	hasMax := len(p) <= 32
	return hasMin && hasMax
}
