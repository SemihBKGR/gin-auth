package auth

import "golang.org/x/crypto/bcrypt"

const defaultCost = 10

type PasswordEncoder interface {
	Encode(password string) (string, error)
	Compare(hash, raw string) bool
}

type BcryptPasswordEncoder struct {
	cost int
}

func (e *BcryptPasswordEncoder) Encode(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), e.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (e *BcryptPasswordEncoder) Compare(hash, raw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
	return err == nil
}

func NewBcryptPasswordEncoder() *BcryptPasswordEncoder {
	return &BcryptPasswordEncoder{
		cost: defaultCost,
	}
}
