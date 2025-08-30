package config

import (
	"asynclab.club/asynx/backend/pkg/entity"
	"asynclab.club/asynx/backend/pkg/util"
)

type ConfigLDAP struct {
	Addr        string `env:"LDAP_ADDR,required"`
	BindDN      string `env:"LDAP_BIND_DN,required"`
	BindPass    string `env:"LDAP_BIND_PASS,required"`
	BaseDN      string `env:"LDAP_BASE_DN,required"`
	UserBaseDN  string `env:"LDAP_USER_BASE_DN,required"`
	GroupBaseDN string `env:"LDAP_GROUP_BASE_DN,required"`
}

var LdapGidNumber = "10000"

var UserObjectClasses = []string{"posixAccount", "inetOrgPerson", "organizationalPerson", "person"}
var GroupObjectClasses = []string{"posixGroup"}

var UserObjectFilter = util.BuildObjectClassCondition(UserObjectClasses)
var GroupObjectFilter = util.BuildObjectClassCondition(GroupObjectClasses)

var UserAttributes []string
var GroupAttributes []string

func init() {
	var err error
	UserAttributes, err = util.GetAttributeKeys[entity.User]()
	if err != nil {
		panic(err)
	}

	GroupAttributes, err = util.GetAttributeKeys[entity.Group]()
	if err != nil {
		panic(err)
	}
}
