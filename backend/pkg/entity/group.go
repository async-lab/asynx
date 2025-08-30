package entity

type Group struct {
	Cn        string   `ldap:"cn,dnAttr:cn,idx:2" json:"cn"`
	Ou        string   `ldap:"ou,dnAttr:ou,idx:1,transient" json:"ou"`
	GidNumber string   `ldap:"gidNumber" json:"gidNumber"`
	MemberUid []string `ldap:"memberUid" json:"memberUid"`
}
