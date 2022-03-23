package auth

import "gin-auth/persist"

type LoginService interface {
	Login(username, password string) (*persist.User, bool)
}

type DefaultLoginService struct {
	userRepo    persist.UserRepository
	passEncoder PasswordEncoder
}

func (s *DefaultLoginService) Login(username, password string) (*persist.User, bool) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil || user == nil {
		return user, false
	}
	return user, s.passEncoder.Compare(user.Password, password)
}

func NewDefaultLoginService(userRepo persist.UserRepository, passEncoder PasswordEncoder) LoginService {
	return &DefaultLoginService{
		userRepo:    userRepo,
		passEncoder: passEncoder,
	}
}
