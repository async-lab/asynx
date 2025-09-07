package controller

import (
	"net/http"

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
	g.GET("", security.GuardMiddleware(security.RoleDefault), gggin.ToGinHandler(ctl.HandleListProfiles))
	g.POST("", security.GuardMiddleware(security.RoleAdmin), gggin.ToGinHandler(ctl.HandleRegister))
	g.GET("/:uid", security.GuardMiddleware(security.RoleRestricted), gggin.ToGinHandler(ctl.HandleGetProfile))
	g.DELETE("/:uid", security.GuardMiddleware(security.RoleAdmin), gggin.ToGinHandler(ctl.HandleUnregister))
	g.PUT("/:uid/password", security.GuardMiddleware(security.RoleRestricted), gggin.ToGinHandler(ctl.HandleChangePassword))
	g.PUT("/:uid/category", security.GuardMiddleware(security.RoleAdmin), gggin.ToGinHandler(ctl.HandleModifyCategory))
	g.PUT("/:uid/role", security.GuardMiddleware(security.RoleAdmin), gggin.ToGinHandler(ctl.HandleModifyRole))

	// Deprecated
	g.GET("/:uid/category", security.GuardMiddleware(security.RoleRestricted), gggin.ToGinHandler(ctl.HandleGetCategory))
	g.GET("/:uid/role", security.GuardMiddleware(security.RoleRestricted), gggin.ToGinHandler(ctl.HandleGetRole))
	return ctl
}

// @Summary      获取用户列表
// @Description  获取所有用户列表信息（包含角色和类别）。需要 ADMIN 角色权限才能查看所有用户，DEFAULT 用户只能查看自己组织单元的用户。
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object} object{data=[]service.UserProfile} "成功返回用户列表"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users [get]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleListProfiles(c *gin.Context) (*gggin.Response[[]*service.UserProfile], *gggin.HttpError) {
	guard, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}

	profiles, err := ctl.serviceManager.ListProfiles(guard)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
	}

	return gggin.NewResponse(profiles), nil
}

// @Summary      获取用户信息
// @Description  根据用户ID获取用户详细信息（包含角色和类别）。需要 RESTRICTED 或更高权限。ADMIN 用户可以查看所有用户信息，DEFAULT 用户只能查看自己组织单元的用户信息，RESTRICTED 用户只能查看自己的信息。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID，使用 'me' 可获取当前用户信息"
// @Success      200  {object} object{data=service.UserProfile} "成功返回用户信息"
// @Failure      400  {object} object{data=string} "请求参数错误"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      404  {object} object{data=string} "用户不存在"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users/{uid} [get]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleGetProfile(c *gin.Context) (*gggin.Response[*service.UserProfile], *gggin.HttpError) {
	guard, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}

	uid := c.Param("uid")
	if uid == "me" {
		uid = guard.Uid
	}

	profile, err := ctl.serviceManager.GetProfile(guard, uid)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
	}

	return gggin.NewResponse(profile), nil
}

type RequestChangePassword struct {
	Password string `json:"password" binding:"required"`
}

// @Summary      修改密码
// @Description  修改指定用户的密码。需要 RESTRICTED 或更高权限。ADMIN 用户可以修改任何用户密码，其他用户只能修改自己的密码。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID，使用 'me' 可修改当前用户密码"
// @Param        body  body      RequestChangePassword  true  "修改密码请求"
// @Success      200  {object} object{data=string} "成功修改密码，返回 'ok'"
// @Failure      400  {object} object{data=string} "请求参数错误或密码不符合要求"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      404  {object} object{data=string} "用户不存在"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users/{uid}/password [put]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleChangePassword(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	guard, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}

	uid := c.Param("uid")
	if uid == "me" {
		uid = guard.Uid
	}
	if guard.Role != security.RoleAdmin && guard.Uid != uid {
		return nil, gggin.NewHttpError(http.StatusForbidden, "权限不足")
	}

	req, err := gggin.ShouldBindJSON[RequestChangePassword](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	err = ctl.serviceManager.ChangePassword(uid, req.Password)
	if err != nil {
		return nil, service.MapErrorToHttp(err)

	}

	return gggin.Ok, nil
}

type RequestModifyCategory struct {
	Category string `json:"category" binding:"required"`
}

// @Summary      更改账号类型
// @Description  修改指定用户的账号类型。需要 ADMIN 角色权限。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID，不能使用 'me'"
// @Param        body  body      RequestModifyCategory  true  "修改账号类型请求\nsystem|member|external"
// @Success      200  {object} object{data=string} "成功修改账号类型，返回 'ok'"
// @Failure      400  {object} object{data=string} "请求参数错误"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      404  {object} object{data=string} "用户不存在"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users/{uid}/category [put]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleModifyCategory(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	guard, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}

	uid := c.Param("uid")
	if uid == "me" || uid == guard.Uid {
		return nil, ErrHttpForceForbidden
	}

	req, err := gggin.ShouldBindJSON[RequestModifyCategory](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	err = ctl.serviceManager.ModifyCategory(uid, req.Category)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
	}

	return gggin.Ok, nil
}

type RequestModifyRole struct {
	Role string `json:"role" binding:"required"`
}

// @Summary      更改账号角色
// @Description  修改指定用户的账号角色。需要 ADMIN 角色权限。非SYSTEM用户必须用学号作为用户名。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true "用户ID，不能使用 'me'"
// @Param        body  body      RequestModifyRole  true  "修改账号角色请求\nadmin|default|restricted"
// @Success      200  {object} object{data=string} "成功修改账号角色，返回 'ok'"
// @Failure      400  {object} object{data=string} "请求参数错误"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      404  {object} object{data=string} "用户不存在"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users/{uid}/role [put]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleModifyRole(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	guard, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}

	uid := c.Param("uid")
	if uid == "me" || uid == guard.Uid {
		return nil, ErrHttpForceForbidden
	}

	req, err := gggin.ShouldBindJSON[RequestModifyRole](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	err = ctl.serviceManager.GrantRoleByUidAndRoleName(uid, req.Role)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
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
// @Description  创建新用户账号。需要 ADMIN 角色权限。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      RequestRegister  true  "注册用户请求"
// @Success      200  {object} object{data=string} "成功注册用户，返回 'ok'"
// @Failure      400  {object} object{data=string} "请求参数错误"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      409  {object} object{data=string} "用户已存在"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users [post]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleRegister(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	_, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}
	req, err := gggin.ShouldBindJSON[RequestRegister](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	err = ctl.serviceManager.Register(req.Username, req.SurName, req.GivenName, req.Mail, req.Category, req.Role)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
	}

	return gggin.Ok, nil
}

// @Summary      删除用户
// @Description  删除指定用户账号。需要 ADMIN 角色权限。不允许删除当前登录用户。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID，不能使用 'me'"
// @Success      200  {object} object{data=string} "成功删除用户，返回 'ok'"
// @Failure      400  {object} object{data=string} "请求参数错误"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      404  {object} object{data=string} "用户不存在"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users/{uid} [delete]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleUnregister(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	guard, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}

	uid := c.Param("uid")
	if uid == "me" || uid == guard.Uid {
		return nil, ErrHttpForceForbidden
	}

	err := ctl.serviceManager.Unregister(uid)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
	}

	return gggin.Ok, nil
}

// ----------------------------------------------------------------------------------------------------------------------

// @Deprecated
// @Summary      获取账号角色
// @Description  获取指定用户的账号角色（admin|default|restricted）。需要 RESTRICTED 或更高权限。ADMIN 用户可以查看所有用户信息，DEFAULT 用户只能查看自己组织单元的用户信息，RESTRICTED 用户只能查看自己的信息。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID，使用 'me' 可获取当前用户角色"
// @Success      200  {object} object{data=string} "成功返回账号角色: admin|default|restricted"
// @Failure      400  {object} object{data=string} "请求参数错误"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users/{uid}/role [get]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleGetRole(c *gin.Context) (*gggin.Response[security.Role], *gggin.HttpError) {
	guard, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}

	uid := c.Param("uid")
	if uid == "me" {
		uid = guard.Uid
	}

	user, err := ctl.serviceManager.GetUserWithGuard(guard, uid)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
	}

	resRole, err := ctl.serviceManager.GetRole(user)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
	}

	return gggin.NewResponse(resRole), nil
}

// @Deprecated
// @Summary      获取账号类型
// @Description  获取指定用户的账号类型（system|member|external）。需要 RESTRICTED 或更高权限。ADMIN 用户可以查看所有用户信息，DEFAULT 用户只能查看自己组织单元的用户信息，RESTRICTED 用户只能查看自己的信息。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        uid   path      string  true  "用户ID，使用 'me' 可获取当前用户类型"
// @Success      200  {object} object{data=string} "成功返回账号类型: system|member|external"
// @Failure      400  {object} object{data=string} "请求参数错误"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      404  {object} object{data=string} "用户不存在"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users/{uid}/category [get]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleGetCategory(c *gin.Context) (*gggin.Response[security.OuUser], *gggin.HttpError) {
	guard, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}

	uid := c.Param("uid")
	if uid == "me" {
		uid = guard.Uid
	}

	user, err := ctl.serviceManager.GetUserWithGuard(guard, uid)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
	}

	category, err := security.GetOuUserFromName(user.Ou)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
	}

	return gggin.NewResponse(category), nil
}
