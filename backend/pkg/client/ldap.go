package client

import (
	"fmt"

	"asynclab.club/asynx/backend/pkg/config"
	"github.com/dsx137/gg-kit/pkg/ggkit"
	"github.com/go-ldap/ldap/v3"
	"github.com/sirupsen/logrus"
)

func validateConfig(cfg *config.ConfigLDAP) error {
	if cfg.Addr == "" {
		return fmt.Errorf("LDAP address is required")
	}
	if cfg.BindDN == "" {
		return fmt.Errorf("LDAP bind DN is required")
	}
	if cfg.BaseDN == "" {
		return fmt.Errorf("LDAP base DN is required")
	}
	return nil
}

// ---------------------------------------------------------------------------------------

type LdapClient struct {
	connPool *ggkit.ReusePool[ldap.Conn]
	cfg      *config.ConfigLDAP
}

func NewLdapClient(cfg *config.ConfigLDAP) (*LdapClient, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	pool, err := ggkit.NewReusePool(
		func() (*ldap.Conn, error) {
			conn, err := ldap.DialURL(cfg.Addr)
			if err != nil {
				logrus.Errorf("failed to dial LDAP server: %v", err)
				return nil, fmt.Errorf("failed to dial LDAP server")
			}

			if err := conn.Bind(cfg.BindDN, cfg.BindPass); err != nil {
				conn.Close()
				return nil, fmt.Errorf("failed to bind with admin credentials: %w", err)
			}
			return conn, nil
		},
		func(conn *ldap.Conn) bool {
			if conn == nil || conn.IsClosing() {
				return false
			}
			_, err := conn.WhoAmI(nil)
			return err == nil
		},
		func(conn *ldap.Conn) error {
			if conn != nil {
				return conn.Close()
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &LdapClient{connPool: pool, cfg: cfg}, nil
}

func (c *LdapClient) withConnection(fn func(*ldap.Conn) error) error {
	conn, err := c.connPool.Get()
	if err != nil {
		return fmt.Errorf("failed to get connection from pool: %w", err)
	}
	defer c.connPool.Put(conn)
	return fn(conn)
}

func (c *LdapClient) BuildDn(rdn string) string {
	if rdn != "" {
		return fmt.Sprintf("%s,%s", rdn, c.cfg.BaseDN)
	}
	return c.cfg.BaseDN
}

func (c *LdapClient) GetUserBaseDn() string  { return c.cfg.UserBaseDN }
func (c *LdapClient) GetGroupBaseDn() string { return c.cfg.GroupBaseDN }
func (c *LdapClient) Close() error           { return c.connPool.Close() }

func (c *LdapClient) Authenticate(dn, password string) (bool, error) {
	if dn == "" || password == "" {
		return false, nil
	}

	authConn, err := ldap.DialURL(c.cfg.Addr)
	if err != nil {
		return false, fmt.Errorf("failed to dial LDAP server: %w", err)
	}
	defer authConn.Close()

	err = authConn.Bind(dn, password)
	if err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			return false, nil
		}
		return false, fmt.Errorf("authentication error: %w", err)
	}

	return true, nil
}

func (c *LdapClient) Search(baseDN string, filter string, attributes []string) (*ldap.SearchResult, error) {
	var result *ldap.SearchResult
	err := c.withConnection(func(conn *ldap.Conn) error {
		if len(attributes) == 0 {
			attributes = []string{"dn", "cn", "mail", "displayName"}
		}

		searchRequest := ldap.NewSearchRequest(
			baseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0, 0, false,
			filter,
			attributes,
			nil,
		)

		var err error
		result, err = conn.Search(searchRequest)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}
		return nil
	})
	return result, err
}

func (c *LdapClient) Add(dn string, objectClass []string, attributes map[string][]string) error {
	return c.withConnection(func(conn *ldap.Conn) error {
		addRequest := ldap.NewAddRequest(dn, nil)
		addRequest.Attribute("objectClass", objectClass)
		for attr, values := range attributes {
			addRequest.Attribute(attr, values)
		}
		return conn.Add(addRequest)
	})
}

func (c *LdapClient) Modify(dn string, addAttrs, delAttrs, replaceAttrs map[string][]string) error {
	return c.withConnection(func(conn *ldap.Conn) error {
		modifyReq := ldap.NewModifyRequest(dn, nil)
		for attr, values := range addAttrs {
			modifyReq.Add(attr, values)
		}
		for attr, values := range delAttrs {
			modifyReq.Delete(attr, values)
		}
		for attr, values := range replaceAttrs {
			modifyReq.Replace(attr, values)
		}
		return conn.Modify(modifyReq)
	})
}

func (c *LdapClient) Delete(dn string) error {
	return c.withConnection(func(conn *ldap.Conn) error {
		delRequest := ldap.NewDelRequest(dn, nil)
		return conn.Del(delRequest)
	})
}

func (c *LdapClient) ModifyDn(dn, newRDN, newSuperior string) error {
	return c.withConnection(func(conn *ldap.Conn) error {
		ModifyDnReq := ldap.NewModifyDNRequest(dn, newRDN, true, newSuperior)
		return conn.ModifyDN(ModifyDnReq)
	})
}

func (c *LdapClient) ModifyPassword(dn, newPassword string) error {
	return c.withConnection(func(conn *ldap.Conn) error {
		passwdReq := ldap.NewPasswordModifyRequest(dn, "", newPassword)
		_, err := conn.PasswordModify(passwdReq)
		return err
	})
}
