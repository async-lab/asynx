package security

import (
	"github.com/dsx137/gg-gin/pkg/gggin"
	"github.com/gin-gonic/gin"
)

func Guard(c *gin.Context, role Role) (string, Role, *gggin.HttpError) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", RoleAnonymous, gggin.NewHttpError(401, "Authorization header is missing")
	}
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", RoleAnonymous, gggin.NewHttpError(401, "Authorization header format must be Bearer {token}")
	}

	tokenString := authHeader[len("Bearer "):]
	claims, err := ParsePaseto(tokenString)
	if err != nil {
		return "", RoleAnonymous, gggin.NewHttpError(401, "Invalid token: "+err.Error())
	}

	if !claims.Role.Support(role) {
		return "", RoleAnonymous, gggin.NewHttpError(403, "Insufficient permissions")
	}

	return claims.Uid, claims.Role, nil
}
