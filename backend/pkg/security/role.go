package security

import (
	"fmt"
	"strings"

	"asynclab.club/asynx/backend/pkg/entity"
)

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleDefault    Role = "default"
	RoleRestricted Role = "restricted"
	RoleAnonymous  Role = "anonymous"
)

func AllRoles() []Role { return []Role{RoleAdmin, RoleDefault, RoleRestricted, RoleAnonymous} }

func (r Role) String() string { return string(r) }

func (r Role) Support(other Role) bool {
	switch r {
	case RoleAdmin:
		return other == RoleAdmin || other == RoleDefault || other == RoleRestricted
	case RoleDefault:
		return other == RoleDefault || other == RoleRestricted
	case RoleRestricted:
		return other == RoleRestricted
	default:
		return false
	}
}

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleDefault, RoleRestricted, RoleAnonymous:
		return true
	default:
		return false
	}
}

func (r Role) IsHigherThan(other Role) bool {
	if !r.IsValid() || !other.IsValid() {
		return false
	}
	switch r {
	case RoleAdmin:
		return other != RoleAdmin
	case RoleDefault:
		return other == RoleRestricted || other == RoleAnonymous
	case RoleRestricted:
		return other == RoleAnonymous
	default:
		return false
	}
}

func GetRoleFromName(authority string) (Role, error) {
	switch strings.ToLower(authority) {
	case "admin":
		return RoleAdmin, nil
	case "default":
		return RoleDefault, nil
	case "restricted":
		return RoleRestricted, nil
	default:
		return RoleAnonymous, fmt.Errorf("unknown role: %s", authority)
	}
}

func GetRoleFromLdapGroup(ldapGroup *entity.Group) (Role, error) {
	return GetRoleFromName(ldapGroup.Cn)
}

func GetRoleFromLdapGroups(ldapGroups []*entity.Group) (Role, error) {
	highestRole, valid := RoleAnonymous, false
	for _, ldapGroup := range ldapGroups {
		if role, err := GetRoleFromLdapGroup(ldapGroup); err == nil {
			if role.IsHigherThan(highestRole) {
				highestRole = role
				valid = true
			}
		}
	}
	if !valid {
		return highestRole, fmt.Errorf("no valid role found")
	}
	return highestRole, nil
}
