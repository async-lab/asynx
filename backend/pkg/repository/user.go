package repository

import (
	"fmt"

	"asynclab.club/asynx/backend/pkg/client"
	"asynclab.club/asynx/backend/pkg/config"
	"asynclab.club/asynx/backend/pkg/entity"
	"asynclab.club/asynx/backend/pkg/transfer"
)

type RepositoryUser struct {
	client *client.LdapClient
}

func NewRepositoryUser(client *client.LdapClient) *RepositoryUser {
	return &RepositoryUser{
		client: client,
	}
}

func (r *RepositoryUser) GetUserBaseDn() string {
	return r.client.GetUserBaseDn()
}

func (r *RepositoryUser) BuildDn(user *entity.User) string {
	return fmt.Sprintf("cn=%s,ou=%s,%s", user.Cn, user.Ou, r.GetUserBaseDn())
}

func (r *RepositoryUser) Authenticate(uid, password string) (bool, error) {
	user, err := r.FindByUid(uid)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return r.client.Authenticate(r.BuildDn(user), password)
}

func (r *RepositoryUser) find(rdn string, filter string) ([]*entity.User, error) {
	baseDN := r.GetUserBaseDn()
	if rdn != "" {
		baseDN = fmt.Sprintf("%s,%s", rdn, baseDN)
	}

	result, err := r.client.Search(baseDN, fmt.Sprintf("(&(%s)(%s))", config.UserObjectFilter, filter), config.UserAttributes)
	if err != nil {
		return nil, err
	}
	users := make([]*entity.User, 0, len(result.Entries))
	for _, entry := range result.Entries {
		user, err := transfer.ParseFromLdap[entity.User](entry)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *RepositoryUser) FindByUid(uid string) (user *entity.User, err error) {
	users, err := r.find("", fmt.Sprintf("uid=%s", uid))
	if len(users) != 0 {
		user = users[0]
	}
	return
}

func (r *RepositoryUser) FindByOuAndUid(ou string, uid string) (user *entity.User, err error) {
	users, err := r.find(fmt.Sprintf("ou=%s", ou), fmt.Sprintf("(uid=%s)", uid))
	if len(users) != 0 {
		user = users[0]
	}
	return
}

func (r *RepositoryUser) FindAll() ([]*entity.User, error) {
	return r.find("", "objectClass=*")
}

func (r *RepositoryUser) FindAllByOu(ou string) ([]*entity.User, error) {
	return r.find(fmt.Sprintf("ou=%s", ou), "objectClass=*")
}

func (r *RepositoryUser) Create(user *entity.User) error {
	attributes, err := transfer.ParseToLdapAttributes(user)
	if err != nil {
		return err
	}

	return r.client.Add(r.BuildDn(user), config.UserObjectClasses, attributes)
}

func (r *RepositoryUser) Modify(user *entity.User) error {
	attributes, err := transfer.ParseToLdapAttributes(user)
	if err != nil {
		return err
	}

	return r.client.Modify(r.BuildDn(user), nil, nil, attributes)
}

func (r *RepositoryUser) ModifyDn(user *entity.User, newRDN, newSuperior string) error {
	return r.client.ModifyDn(r.BuildDn(user), newRDN, newSuperior)
}

func (r *RepositoryUser) ModifyPassword(user *entity.User, newPassword string) error {
	return r.client.ModifyPassword(r.BuildDn(user), newPassword)
}

func (r *RepositoryUser) Delete(user *entity.User) error {
	return r.client.Delete(r.BuildDn(user))
}
