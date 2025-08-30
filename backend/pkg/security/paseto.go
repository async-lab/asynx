package security

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"asynclab.club/asynx/backend/pkg/config"
)

type PasetoClaims struct {
	Uid  string `json:"uid"`
	Role Role   `json:"role"`
}

func GeneratePaseto(uid string, role Role) (string, error) {
	token := paseto.NewToken()

	token.SetSubject(uid)
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(time.Hour * 24))
	token.SetIssuer("asynx")

	token.Set("claims", &PasetoClaims{
		Uid:  uid,
		Role: role,
	})

	return token.V4Encrypt(config.PasetoKey, nil), nil
}

func ParsePaseto(tokenString string) (*PasetoClaims, error) {
	parser := paseto.NewParser()

	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.IssuedBy("asynx"))

	parsedToken, err := parser.ParseV4Local(config.PasetoKey, tokenString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	var claims PasetoClaims
	err = parsedToken.Get("claims", &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &claims, nil
}
