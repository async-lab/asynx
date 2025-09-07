package entity

type User struct {
	Uid           string `ldap:"uid" json:"uid"`
	Cn            string `ldap:"cn,dnAttr:cn,idx:2" json:"cn"`
	Ou            string `ldap:"ou,dnAttr:ou,idx:1" json:"ou"`
	Sn            string `ldap:"sn" json:"sn"`
	GivenName     string `ldap:"givenName" json:"givenName"`
	GidNumber     string `ldap:"gidNumber" json:"gidNumber"`
	UidNumber     string `ldap:"uidNumber" json:"uidNumber"`
	HomeDirectory string `ldap:"homeDirectory" json:"homeDirectory"`
	Mail          string `ldap:"mail" json:"mail"`
	UserPassword  string `ldap:"userPassword" json:"userPassword"`
}
