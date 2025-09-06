package security

import (
	"github.com/dsx137/gg-gin/pkg/gggin"
	"github.com/gin-gonic/gin"
)

type GuardResult struct {
	Uid  string
	Role Role
}

func Guard(c *gin.Context, role Role) (string, Role, *gggin.HttpError) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", RoleAnonymous, gggin.NewHttpError(401, "缺少Authorization头")
	}
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", RoleAnonymous, gggin.NewHttpError(401, "Authorization头格式必须为: Bearer {token}")
	}

	tokenString := authHeader[len("Bearer "):]
	claims, err := ParsePaseto(tokenString)
	if err != nil {
		return "", RoleAnonymous, gggin.NewHttpError(401, "无效的令牌: "+err.Error())
	}

	if !claims.Role.Support(role) {
		return "", RoleAnonymous, gggin.NewHttpError(403, "权限不足")
	}

	return claims.Uid, claims.Role, nil
}

func GuardMiddleware(role Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, role, err := Guard(c, role)
		if err != nil {
			c.JSON(err.StatusCode, err.Message)
			c.Abort()
			return
		}

		c.Set("guard", &GuardResult{Uid: uid, Role: role})
		c.Next()
	}
}
