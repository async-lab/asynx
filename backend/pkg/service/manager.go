package service

import (
	"fmt"
	"net/http"
	"strconv"

	"asynclab.club/asynx/backend/pkg/client"
	"asynclab.club/asynx/backend/pkg/config"
	"asynclab.club/asynx/backend/pkg/entity"
	"asynclab.club/asynx/backend/pkg/security"
	"asynclab.club/asynx/backend/pkg/util"
	"github.com/dsx137/gg-gin/pkg/gggin"
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

func (s *ServiceManager) Authenticate(username string, password string) (string, error) {
	ok, err := s.serviceUser.Authenticate(username, password)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", gggin.NewHttpError(http.StatusUnauthorized, "无效的凭证")
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
		return gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	if ou != security.OuUserSystem {
		if err := security.ValidateMemberUsernameLegality(username); err != nil {
			return gggin.NewHttpError(http.StatusBadRequest, err.Error())
		}
	}

	role, err := security.GetRoleFromName(roleName)
	if err != nil {
		return gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	err = security.ValidateEmailFormat(mail)
	if err != nil {
		return gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	ok, err := s.CheckUserExists(username)
	if err != nil {
		return err
	}
	if ok {
		return gggin.NewHttpError(http.StatusBadRequest, "用户已存在")
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

	err = s.serviceUser.Create(user)
	if err != nil {
		return err
	}
	err = s.serviceGroup.GrantRole(user, role)
	if err != nil {
		s.Unregister(user)
		return err
	}
	err = s.emailClient.SendMail(
		user.Mail,
		"异步实验室",
		struct{ Surname, GivenName, Username, Password string }{
			Surname:   user.Sn,
			GivenName: user.GivenName,
			Username:  user.Uid,
			Password:  user.UserPassword,
		},
	)
	if err != nil {
		if err := s.Unregister(user); err != nil {
			logrus.Warningf("Failed to rollback when register: %v", err)
		}
		return err
	}
	return nil
}

func (s *ServiceManager) Unregister(user *entity.User) error {
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

func (s *ServiceManager) GetRoleByUid(uid string) (security.Role, error) {
	return s.serviceGroup.GetRoleByUid(uid)
}

func (s *ServiceManager) FindUserByUid(uid string) (*entity.User, error) {
	user, err := s.serviceUser.FindByUid(uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gggin.NewHttpError(http.StatusNotFound, "用户未找到")
	}
	return user, nil
}

func (s *ServiceManager) FindUserByOuAndUid(ou security.OuUser, uid string) (*entity.User, error) {
	user, err := s.serviceUser.FindByOuAndUid(ou, uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gggin.NewHttpError(http.StatusNotFound, "用户未找到")
	}
	return user, nil
}

func (s *ServiceManager) FindAllUsers() ([]*entity.User, error) {
	return s.serviceUser.FindAll()
}

func (s *ServiceManager) FindAllUsersByOu(ou security.OuGroup) ([]*entity.User, error) {
	return s.serviceUser.FindAllByOu(ou)
}

func (s *ServiceManager) ModifyPassword(user *entity.User, newPassword string) error {
	return s.serviceUser.ModifyPassword(user, newPassword)
}

func (s *ServiceManager) ModifyCategory(user *entity.User, ou security.OuUser) error {
	return s.serviceUser.ModifyCategory(user, ou)
}

func (s *ServiceManager) GetRole(user *entity.User) (security.Role, error) {
	return s.serviceGroup.GetRole(user)
}

func (s *ServiceManager) GrantRoleByUidAndRoleName(uid string, roleName string) error {
	user, err := s.FindUserByUid(uid)
	if err != nil {
		return err
	}
	if user == nil {
		return gggin.NewHttpError(http.StatusNotFound, "用户未找到")
	}

	role, err := security.GetRoleFromName(roleName)
	if err != nil {
		return gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	err = s.serviceGroup.GrantRole(user, role)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceManager) GenerateNextUidNumber() (string, error) {
	users, err := s.FindAllUsers()
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
	users, err := s.FindAllUsers()
	if err != nil {
		return false, err
	}

	for _, user := range users {
		if user.Uid == username {
			return true, nil
		}
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
		user, err = s.FindUserByUid(uid)
	case security.RoleDefault:
		user, err = s.FindUserByOuAndUid(security.OuUserMember, uid)
	default:
		if authUid != uid {
			return nil, nil
		}

		user, err = s.FindUserByUid(uid)
	}

	return user, err
}
