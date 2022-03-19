package auth

import "gin-auth/persist"

type LoginService interface {
	login(username, password string) bool
}

type DefaultLoginService struct {
	userRepo    persist.UserRepository
	passEncoder *PasswordEncoder
}

func (s DefaultLoginService) login(username, password string) bool {
	user := s.userRepo.FindByUsername(username)
	if user == nil {
		return false
	}
	return s.passEncoder.Compare(user.Password, password)
}

func NewDefaultLoginService(userRepo persist.UserRepository, passEncoder *PasswordEncoder) *DefaultLoginService {
	return &DefaultLoginService{
		userRepo:    userRepo,
		passEncoder: passEncoder,
	}
}
