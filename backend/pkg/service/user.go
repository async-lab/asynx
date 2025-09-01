package service

import (
	"fmt"

	"asynclab.club/asynx/backend/pkg/entity"
	"asynclab.club/asynx/backend/pkg/repository"
	"asynclab.club/asynx/backend/pkg/security"
)

type ServiceUser struct {
	repositoryUser *repository.RepositoryUser
}

func NewServiceUser(repo *repository.RepositoryUser) *ServiceUser {
	return &ServiceUser{repositoryUser: repo}
}

func (s *ServiceUser) Authenticate(uid, password string) (bool, error) {
	return s.repositoryUser.Authenticate(uid, password)
}

func (s *ServiceUser) FindByUid(uid string) (*entity.User, error) {
	user, err := s.repositoryUser.FindByUid(uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *ServiceUser) FindByOuAndUid(ou security.OuUser, uid string) (*entity.User, error) {
	user, err := s.repositoryUser.FindByOuAndUid(ou.String(), uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *ServiceUser) FindAll() ([]*entity.User, error) {
	return s.repositoryUser.FindAll()
}

func (s *ServiceUser) FindAllByOu(ou security.OuUser) ([]*entity.User, error) {
	return s.repositoryUser.FindAllByOu(ou.String())
}

func (s *ServiceUser) Create(user *entity.User) error {
	return s.repositoryUser.Create(user)
}

func (s *ServiceUser) Modify(user *entity.User) error {
	return s.repositoryUser.Modify(user)
}

func (s *ServiceUser) ModifyPassword(user *entity.User, newPassword string) error {
	return s.repositoryUser.ModifyPassword(user, newPassword)
}

func (s *ServiceUser) ModifyOu(user *entity.User, ou security.OuUser) error {
	return s.repositoryUser.ModifyDn(user, fmt.Sprintf("cn=%s", user.Cn), fmt.Sprintf("ou=%s,%s", ou.String(), s.repositoryUser.GetUserBaseDn()))
}

func (s *ServiceUser) Delete(user *entity.User) error {
	return s.repositoryUser.Delete(user)
}
