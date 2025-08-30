package controller

import (
	"errors"
	"net/http"

	"asynclab.club/asynx/backend/pkg/entity"
	"asynclab.club/asynx/backend/pkg/security"
	"asynclab.club/asynx/backend/pkg/service"
	"github.com/dsx137/gg-gin/pkg/gggin"
	"github.com/gin-gonic/gin"
)

type ControllerUser struct {
	serviceManager *service.ServiceManager
}

func NewControllerUser(g *gin.RouterGroup, serviceManager *service.ServiceManager) *ControllerUser {
	ctl := &ControllerUser{serviceManager: serviceManager}
	g.GET("", gggin.HandleController(ctl.HandleList))
	g.POST("", gggin.HandleController(ctl.HandleRegister))
	g.GET("/:uid", gggin.HandleController(ctl.HandleGet))
	g.DELETE("/:uid", gggin.HandleController(ctl.HandleUnregister))
	g.PUT("/:uid/password", gggin.HandleController(ctl.HandleChangePassword))
	g.GET("/:uid/category", gggin.HandleController(ctl.HandleGetCategory))
	g.PUT("/:uid/category", gggin.HandleController(ctl.HandleModifyCategory))
	g.GET("/:uid/role", gggin.HandleController(ctl.HandleGetRole))
	g.PUT("/:uid/role", gggin.HandleController(ctl.HandleModifyRole))
	return ctl
}

func (ctl *ControllerUser) processGetWithAuthority(authUid string, uid string, role security.Role) (*entity.User, *gggin.HttpError) {
	if uid == "me" {
		uid = authUid
	}

	var (
		user *entity.User
		err  error
	)

	switch role {
	case security.RoleAdmin:
		user, err = ctl.serviceManager.FindUserByUid(uid)
	case security.RoleDefault:
		user, err = ctl.serviceManager.FindUserByOuAndUid(security.OuUserMember, uid)
	default:
		if authUid != uid {
			return nil, nil
		}

		user, err = ctl.serviceManager.FindUserByUid(uid)
	}

	if err != nil {
		var httpErr *gggin.HttpError
		if errors.As(err, &httpErr) {
			return nil, httpErr
		}
		return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	return user, nil
}

// @Summary      获取用户列表
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object} object{data=[]entity.User}
// @Router       /api/users [get]
func (ctl *ControllerUser) HandleList(c *gin.Context) (*gggin.Response[[]*entity.User], *gggin.HttpError) {
	_, role, guardErr := security.Guard(c, security.RoleDefault)
	if guardErr != nil {
		return nil, guardErr
	}

	switch role {
	case security.RoleAdmin:
		users, err := ctl.serviceManager.FindAllUsers()
		if err != nil {
			return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
		}
		return gggin.NewResponse(users), nil
	default:
		users, err := ctl.serviceManager.FindAllUsersByOu(security.OuGroup(security.OuUserMember))
		if err != nil {
			return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
		}
		return gggin.NewResponse(users), nil
	}
}

// @Summary      获取用户信息
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "用户ID"
// @Success      200  {object} object{data=entity.User}
// @Router       /api/users/{id} [get]
func (ctl *ControllerUser) HandleGet(c *gin.Context) (*gggin.Response[*entity.User], *gggin.HttpError) {
	authUid, role, guardErr := security.Guard(c, security.RoleRestricted)
	if guardErr != nil {
		return nil, guardErr
	}

	user, err := ctl.processGetWithAuthority(authUid, c.Param("uid"), role)
	if err != nil {
		return nil, err
	}

	return gggin.NewResponse(user), nil
}

type RequestChangePassword struct {
	Password string `json:"password" binding:"required"`
}

// @Summary      修改密码
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID"
// @Param        body  body      RequestChangePassword  true  "修改密码请求"
// @Success      200  {object} object{data=string} "Success response with 'ok'"
// @Router       /api/users/{id}/password [put]
func (ctl *ControllerUser) HandleChangePassword(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	authUid, role, guardErr := security.Guard(c, security.RoleDefault)
	if guardErr != nil {
		return nil, guardErr
	}

	uid := c.Param("uid")
	if uid == "me" {
		uid = authUid
	}
	if role != security.RoleAdmin && authUid != uid {
		return nil, gggin.NewHttpError(http.StatusForbidden, "Insufficient permissions")
	}

	req, err := gggin.ShouldBindJSON[RequestChangePassword](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}
	user, err := ctl.serviceManager.FindUserByUid(uid)
	if err != nil {
		var httpErr *gggin.HttpError
		if errors.As(err, &httpErr) {
			return nil, httpErr
		}
		return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	if err := security.ValidatePasswordLegality(req.Password); err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}
	if err := security.ValidatePasswordStrength(req.Password); err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}
	if err := ctl.serviceManager.ModifyPassword(user, req.Password); err != nil {
		return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	return gggin.Ok, nil
}

// @Summary      获取账号类型
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID"
// @Success      200  {object} object{data=string} ""system|member|external""
// @Router       /api/users/{id}/category [get]
func (ctl *ControllerUser) HandleGetCategory(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	authUid, role, guardErr := security.Guard(c, security.RoleRestricted)
	if guardErr != nil {
		return nil, guardErr
	}

	user, processErr := ctl.processGetWithAuthority(authUid, c.Param("uid"), role)
	if processErr != nil {
		return nil, processErr
	}

	category, err := security.GetOuUserFromName(user.Ou)
	if err != nil {
		return nil, gggin.NewHttpError(500, err.Error())
	}

	return gggin.NewResponse(category.String()), nil
}

type RequestModifyCategory struct {
	Category string `json:"category" binding:"required"`
}

// @Summary      更改账号类型
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID"
// @Param        body  body      RequestModifyCategory  true  "修改账号类型请求\nsystem|member|external"
// @Success      200  {object} object{data=string} "Success response with 'ok'"
// @Router       /api/users/{id}/category [put]
func (ctl *ControllerUser) HandleModifyCategory(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	authUid, _, guardErr := security.Guard(c, security.RoleAdmin)
	if guardErr != nil {
		return nil, guardErr
	}

	uid := c.Param("uid")
	if uid == "me" {
		uid = authUid
	}

	req, err := gggin.ShouldBindJSON[RequestModifyCategory](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	user, err := ctl.serviceManager.FindUserByUid(uid)
	if err != nil {
		var httpErr *gggin.HttpError
		if errors.As(err, &httpErr) {
			return nil, httpErr
		}
		return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	category, err := security.GetOuUserFromName(req.Category)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	err = ctl.serviceManager.ModifyCategory(user, category)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	return gggin.Ok, nil
}

// @Summary      获取账号角色
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID"
// @Success      200  {object} object{data=string} ""admin|default|restricted""
// @Router       /api/users/{id}/role [get]
func (ctl *ControllerUser) HandleGetRole(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	_, role, guardErr := security.Guard(c, security.RoleRestricted)
	if guardErr != nil {
		return nil, guardErr
	}

	return gggin.NewResponse(role.String()), nil
}

type RequestModifyRole struct {
	Role string `json:"role" binding:"required"`
}

// @Summary      更改账号角色
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID"
// @Param        body  body      RequestModifyRole  true  "修改账号角色请求\nadmin|default|restricted"
// @Success      200  {object} object{data=string} "Success response with 'ok'"
// @Router       /api/users/{id}/role [put]
func (ctl *ControllerUser) HandleModifyRole(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	authUid, _, guardErr := security.Guard(c, security.RoleAdmin)
	if guardErr != nil {
		return nil, guardErr
	}

	uid := c.Param("uid")
	if uid == "me" {
		uid = authUid
	}

	req, err := gggin.ShouldBindJSON[RequestModifyRole](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	err = ctl.serviceManager.GrantRoleByUidAndRoleName(uid, req.Role)
	if err != nil {
		var httpErr *gggin.HttpError
		if errors.As(err, &httpErr) {
			return nil, httpErr
		}
		return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	return gggin.Ok, nil
}

type RequestRegister struct {
	Username  string `json:"username" binding:"required"`
	SurName   string `json:"surName" binding:"required"`
	GivenName string `json:"givenName" binding:"required"`
	Mail      string `json:"mail" binding:"required"`
	Category  string `json:"category" binding:"required"`
	Role      string `json:"role" binding:"required"`
}

// @Summary      注册新用户
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body  body      RequestRegister  true  "注册用户请求"
// @Success      200  {object} object{data=string} "Success response with 'ok'"
// @Router       /api/users [post]
func (ctl *ControllerUser) HandleRegister(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	_, _, guardErr := security.Guard(c, security.RoleAdmin)
	if guardErr != nil {
		return nil, guardErr
	}

	req, err := gggin.ShouldBindJSON[RequestRegister](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	err = ctl.serviceManager.Register(req.Username, req.SurName, req.GivenName, req.Mail, req.Category, req.Role)
	if err != nil {
		var httpErr *gggin.HttpError
		if errors.As(err, &httpErr) {
			return nil, httpErr
		}
		return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	return gggin.Ok, nil
}

// @Summary      删除用户
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID"
// @Success      200  {object} object{data=string} "Success response with 'ok'"
// @Router       /api/users/{id} [delete]
func (ctl *ControllerUser) HandleUnregister(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	authUid, _, guardErr := security.Guard(c, security.RoleAdmin)
	if guardErr != nil {
		return nil, guardErr
	}
	uid := c.Param("uid")
	if uid == "me" || uid == authUid {
		return nil, gggin.NewHttpError(http.StatusForbidden, "WHAT ARE YOU DOING?")
	}
	user, err := ctl.serviceManager.FindUserByUid(uid)
	if err != nil {
		var httpErr *gggin.HttpError
		if errors.As(err, &httpErr) {
			return nil, httpErr
		}
		return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	err = ctl.serviceManager.Unregister(user)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	return gggin.Ok, nil
}
