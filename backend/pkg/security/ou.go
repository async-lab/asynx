package security

import "fmt"

type OuUser string

const (
	OuUserSystem   OuUser = "system"
	OuUserMember   OuUser = "member"
	OuUserExternal OuUser = "external"
	OuUserUnknown  OuUser = "unknown"
)

func (u OuUser) String() string { return string(u) }
func GetOuUserFromName(ouName string) (OuUser, error) {
	switch ouName {
	case "system":
		return OuUserSystem, nil
	case "member":
		return OuUserMember, nil
	case "external":
		return OuUserExternal, nil
	default:
		return OuUserUnknown, fmt.Errorf("unknown ou user: %s", ouName)
	}
}

type OuGroup string

const (
	OuGroupPrimary       OuGroup = "primary"
	OuGroupSupplementary OuGroup = "supplementary"
	OuGroupAdditional    OuGroup = "additional"
	OuGroupUnknown       OuGroup = "unknown"
)

func (g OuGroup) String() string { return string(g) }
func GetOuGroupFromName(ouName string) (OuGroup, error) {
	switch ouName {
	case "primary":
		return OuGroupPrimary, nil
	case "supplementary":
		return OuGroupSupplementary, nil
	case "additional":
		return OuGroupAdditional, nil
	default:
		return OuGroupUnknown, fmt.Errorf("unknown ou group: %s", ouName)
	}
}
