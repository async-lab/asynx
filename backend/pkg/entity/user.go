package entity

type User struct {
	Uid           string `ldap:"uid" json:"username"`
	Cn            string `ldap:"cn,dnAttr:cn,idx:2" json:"-"`
	Ou            string `ldap:"ou,dnAttr:ou,idx:1" json:"-"`
	Sn            string `ldap:"sn" json:"surName"`
	GivenName     string `ldap:"givenName" json:"givenName"`
	GidNumber     string `ldap:"gidNumber" json:"-"`
	UidNumber     string `ldap:"uidNumber" json:"-"`
	HomeDirectory string `ldap:"homeDirectory" json:"-"`
	Mail          string `ldap:"mail" json:"mail"`
	UserPassword  string `ldap:"userPassword" json:"-"`
}
