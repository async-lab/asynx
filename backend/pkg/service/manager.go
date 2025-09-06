package service

import (
	"fmt"
	"strconv"

	"asynclab.club/asynx/backend/pkg/client"
	"asynclab.club/asynx/backend/pkg/config"
	"asynclab.club/asynx/backend/pkg/entity"
	"asynclab.club/asynx/backend/pkg/security"
	"asynclab.club/asynx/backend/pkg/util"
	"github.com/dsx137/gg-kit/pkg/ggkit"
	"github.com/sirupsen/logrus"
)

type ServiceManager struct {
	serviceUser  *ServiceUser
	serviceGroup *ServiceGroup
	emailClient  *client.EmailClient
}

func NewServiceManager(serviceUser *ServiceUser, serviceGroup *ServiceGroup, emailClient *client.EmailClient) *ServiceManager {
	return &ServiceManager{serviceUser: serviceUser, serviceGroup: serviceGroup, emailClient: emailClient}
}

func (s *ServiceManager) Authenticate(username, password string) (string, error) {
	ok, err := s.serviceUser.Authenticate(username, password)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrInvalidCreds
	}

	role, err := s.GetRoleByUid(username)
	if err != nil {
		return "", err
	}

	token, err := security.GeneratePaseto(username, role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *ServiceManager) Register(username, surName, givenName, mail, category, roleName string) error {
	ou, err := security.GetOuUserFromName(category)
	if err != nil {
		return ErrInvalidOu
	}

	if ou != security.OuUserSystem {
		if err := security.ValidateMemberUsernameLegality(username); err != nil {
			return fmt.Errorf("%w: %s", ErrConflict, err.Error())
		}
	}

	role, err := security.GetRoleFromName(roleName)
	if err != nil {
		return ErrInvalidRole
	}

	if err := security.ValidateEmailFormat(mail); err != nil {
		return ErrInvalidEmail
	}

	ok, err := s.CheckUserExists(username)
	if err != nil {
		return err
	}
	if ok {
		return ErrUserExists
	}

	uidNumber, err := s.GenerateNextUidNumber()
	if err != nil {
		return err
	}

	password, err := ggkit.GenerateReadableKey(32, 0)
	if err != nil {
		return err
	}

	user := &entity.User{
		Uid:           username,
		Cn:            username,
		Ou:            ou.String(),
		Sn:            surName,
		GivenName:     givenName,
		GidNumber:     config.LdapGidNumber,
		UidNumber:     uidNumber,
		HomeDirectory: fmt.Sprintf("/home/%s", username),
		Mail:          mail,
		UserPassword:  password,
	}

	if err := s.serviceUser.Create(user); err != nil {
		return err
	}

	if err := s.serviceGroup.GrantRole(user, role); err != nil {
		_ = s.unregister(user) // rollback
		return err
	}

	if err := s.emailClient.SendMail(
		user.Mail,
		"异步实验室",
		struct {
			Surname, GivenName, Username, Password string
		}{
			Surname:   user.Sn,
			GivenName: user.GivenName,
			Username:  user.Uid,
			Password:  user.UserPassword,
		},
	); err != nil {
		_ = s.unregister(user) // rollback
		return err
	}

	return nil
}

func (s *ServiceManager) unregister(user *entity.User) error {
	err := s.serviceUser.Delete(user)
	if err != nil {
		return err
	}

	err = s.serviceGroup.GrantRole(user, security.RoleAnonymous)
	if err != nil {
		logrus.Warnf("User %s deleted, but failed to remove from role group: %v", user.Uid, err)
	}
	return nil
}

func (s *ServiceManager) Unregister(uid string) error {
	user, err := s.serviceUser.FindByUid(uid)
	if err != nil {
		return err
	}

	return s.unregister(user)
}

func (s *ServiceManager) GetRoleByUid(uid string) (security.Role, error) {
	return s.serviceGroup.GetRoleByUid(uid)
}
func (s *ServiceManager) GetRole(user *entity.User) (security.Role, error) {
	return s.serviceGroup.GetRole(user)
}

func (s *ServiceManager) GrantRoleByUidAndRoleName(uid string, roleName string) error {
	user, err := s.serviceUser.FindByUid(uid)
	if err != nil {
		return err
	}

	role, err := security.GetRoleFromName(roleName)
	if err != nil {
		return ErrInvalidRole
	}

	err = s.serviceGroup.GrantRole(user, role)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceManager) GenerateNextUidNumber() (string, error) {
	users, err := s.serviceUser.FindAll()
	if err != nil {
		return "", err
	}

	userUidNumbers := make([]int, len(users))
	for i, user := range users {
		userUidNumbers[i], _ = strconv.Atoi(user.UidNumber)
	}

	uidNumber := util.FindFirstMissingPositive(userUidNumbers)
	return strconv.Itoa(uidNumber), nil
}

func (s *ServiceManager) CheckUserExists(username string) (bool, error) {
	user, err := s.serviceUser.FindByUid(username)
	if err == ErrUserNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if user != nil {
		return true, nil
	}

	return false, nil
}

func (s *ServiceManager) GetUserWithAuthority(authUid string, uid string, role security.Role) (*entity.User, error) {
	if uid == "me" {
		uid = authUid
	}

	var (
		user *entity.User
		err  error
	)

	switch role {
	case security.RoleAdmin:
		user, err = s.serviceUser.FindByUid(uid)
	case security.RoleDefault:
		user, err = s.serviceUser.FindByOuAndUid(security.OuUserMember, uid)
	default:
		if authUid != uid {
			return nil, nil
		}

		user, err = s.serviceUser.FindByUid(uid)
	}

	return user, err
}

func (s *ServiceManager) List(uid string, role security.Role) ([]*entity.User, error) {
	switch role {
	case security.RoleAdmin:
		users, err := s.serviceUser.FindAll()
		if err != nil {
			return nil, err
		}
		return users, nil
	default:
		user, err := s.serviceUser.FindByUid(uid)
		if err != nil {
			return nil, err
		}

		ou, err := security.GetOuUserFromName(user.Ou)
		if err != nil {
			return nil, err
		}

		users, err := s.serviceUser.FindAllByOu(ou)
		if err != nil {
			return nil, err
		}
		return users, nil
	}
}

func (s *ServiceManager) ChangePassword(uid string, password string) error {
	user, err := s.serviceUser.FindByUid(uid)
	if err != nil {
		return err
	}

	if err := security.ValidatePasswordLegality(password); err != nil {
		return ErrIllegalPassword
	}
	if err := security.ValidatePasswordStrength(password); err != nil {
		return ErrWeakPassword
	}
	if err := s.serviceUser.ModifyPassword(user, password); err != nil {
		return err
	}
	return nil
}

func (s *ServiceManager) ModifyCategory(uid string, category string) error {
	user, err := s.serviceUser.FindByUid(uid)
	if err != nil {
		return err
	}

	ou, err := security.GetOuUserFromName(category)
	if err != nil {
		return ErrInvalidOu
	}

	err = s.serviceUser.ModifyOu(user, ou)
	if err != nil {
		return err
	}

	return nil
}
