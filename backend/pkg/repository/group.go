package repository

import (
	"fmt"

	"asynclab.club/asynx/backend/pkg/client"
	"asynclab.club/asynx/backend/pkg/config"
	"asynclab.club/asynx/backend/pkg/entity"
	"asynclab.club/asynx/backend/pkg/transfer"
	"asynclab.club/asynx/backend/pkg/util"
)

var groupAttributes []string

func init() {
	var err error
	groupAttributes, err = util.GetAttributeKeys[entity.Group]()
	if err != nil {
		panic(err)
	}
}

type RepositoryGroup struct {
	client *client.LdapClient
}

func NewRepositoryGroup(client *client.LdapClient) *RepositoryGroup {
	return &RepositoryGroup{
		client: client,
	}
}

func (r *RepositoryGroup) GetGroupBaseDn() string {
	return r.client.GetGroupBaseDn()
}

func (r *RepositoryGroup) BuildDn(group *entity.Group) string {
	return fmt.Sprintf("cn=%s,ou=%s,%s", group.Cn, group.Ou, r.GetGroupBaseDn())
}

func (r *RepositoryGroup) find(rdn string, filter string) ([]*entity.Group, error) {
	baseDN := r.GetGroupBaseDn()
	if rdn != "" {
		baseDN = fmt.Sprintf("%s,%s", rdn, baseDN)
	}

	result, err := r.client.Search(baseDN, fmt.Sprintf("(&(%s)(%s))", config.GroupObjectFilter, filter), groupAttributes)
	if err != nil {
		return nil, err
	}
	groups := make([]*entity.Group, 0, len(result.Entries))
	for _, entry := range result.Entries {
		group, err := transfer.ParseFromLdap[entity.Group](entry)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (r *RepositoryGroup) FindByOuAndCn(ou string, cn string) (group *entity.Group, err error) {
	groups, err := r.find(fmt.Sprintf("ou=%s", ou), fmt.Sprintf("cn=%s", cn))
	if len(groups) != 0 {
		group = groups[0]
	}
	return
}

func (r *RepositoryGroup) FindAll() ([]*entity.Group, error) {
	return r.find("", "objectClass=*")
}

func (r *RepositoryGroup) FindAllByOuAndMemberUid(ou string, uid string) ([]*entity.Group, error) {
	return r.find(fmt.Sprintf("ou=%s", ou), fmt.Sprintf("memberUid=%s", uid))
}

func (r *RepositoryGroup) Modify(dn string, addAttrs map[string][]string, delAttrs map[string][]string, replaceAttrs map[string][]string) error {
	return r.client.Modify(dn, addAttrs, delAttrs, replaceAttrs)
}
