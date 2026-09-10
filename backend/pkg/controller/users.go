package controller

import (
	"encoding/csv"
	"io"
	"net/http"
	"strings"

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
	g.PATCH("/category", security.GuardMiddleware(security.RoleAdmin), gggin.ToGinHandler(ctl.HandleBatchModifyCategory))
	g.PATCH("/role", security.GuardMiddleware(security.RoleAdmin), gggin.ToGinHandler(ctl.HandleBatchModifyRole))

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
	if uid == "me" {
		uid = guard.Uid
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
// @Description  修改指定用户的账号角色。需要 ADMIN 角色权限。非SYSTEM用户必须用学号作为用户名。不允许操作当前登录用户。
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
// @Description  创建新用户账号。支持两种方式：1) JSON格式单个注册 2) CSV文件批量注册。需要 ADMIN 角色权限。
// @Description
// @Description  **单个注册 (application/json)**
// @Description  发送JSON格式的单个用户数据
// @Description
// @Description  **批量注册 (multipart/form-data 或 text/csv)**
// @Description  上传CSV文件，CSV格式为：username,surName,givenName,mail,category,role
// @Description  示例：user001,张,三,zhangsan@example.com,member,default
// @Tags         users
// @Accept       json,multipart/form-data,text/csv
// @Produce      json
// @Param        body  body      RequestRegister  false  "单个注册请求 (Content-Type: application/json)"
// @Param        file  formData  file  false  "批量注册CSV文件 (Content-Type: multipart/form-data)"
// @Success      200  {object} object{data=string} "单个注册成功，返回 'ok'"
// @Success      200  {object} object{data=BatchRegisterResult} "批量注册全部成功"
// @Success      207  {object} object{data=BatchRegisterResult} "批量注册部分失败"
// @Failure      400  {object} object{data=string} "请求参数错误"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      409  {object} object{data=string} "用户已存在"
// @Failure      415  {object} object{data=string} "不支持的媒体类型"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users [post]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleRegister(c *gin.Context) (*gggin.Response[any], *gggin.HttpError) {
	_, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}

	contentType := c.ContentType()

	// 根据 Content-Type 区分单个注册还是批量注册
	if contentType == "application/json" {
		// 单个用户注册
		return ctl.handleSingleRegister(c)
	} else if strings.HasPrefix(contentType, "multipart/form-data") || contentType == "text/csv" {
		// CSV 批量注册
		return ctl.handleBatchRegisterFromCSV(c)
	}

	return nil, gggin.NewHttpError(http.StatusUnsupportedMediaType, "不支持的 Content-Type，请使用 application/json (单个注册) 或 multipart/form-data (批量注册)")
}

// handleSingleRegister 处理单个用户注册（JSON格式）
func (ctl *ControllerUser) handleSingleRegister(c *gin.Context) (*gggin.Response[any], *gggin.HttpError) {
	req, err := gggin.ShouldBindJSON[RequestRegister](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	err = ctl.serviceManager.Register(req.Username, req.SurName, req.GivenName, req.Mail, req.Category, req.Role)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
	}

	return gggin.NewResponse[any]("ok"), nil
}

// handleBatchRegisterFromCSV 处理CSV文件批量注册
func (ctl *ControllerUser) handleBatchRegisterFromCSV(c *gin.Context) (*gggin.Response[any], *gggin.HttpError) {
	file, err := c.FormFile("file")
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, "未找到上传的文件，请使用 'file' 字段上传CSV文件")
	}

	src, err := file.Open()
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusInternalServerError, "无法读取上传的文件")
	}
	defer src.Close()

	reader := csv.NewReader(src)

	header, err := reader.Read()
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, "CSV文件格式错误：无法读取表头")
	}

	// 验证表头格式
	expectedHeaders := []string{"username", "surName", "givenName", "mail", "category", "role"}
	hasHeader := false
	if len(header) == len(expectedHeaders) {
		// 检查前3个关键字段来判断是否有表头，避免误判
		h0 := strings.ToLower(header[0])
		h3 := strings.ToLower(header[3])
		h4 := strings.ToLower(header[4])
		if h0 == "username" && h3 == "mail" && h4 == "category" {
			hasHeader = true
		}
	}

	var result BatchRegisterResult
	rowNum := 1

	if !hasHeader {
		if err := ctl.processCSVRow(header, &result, rowNum); err != nil {
			return nil, err
		}
		rowNum++
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Failed++
			result.Total++
			result.Failures = append(result.Failures, BatchFailureInfo{
				Row:      rowNum,
				Username: "",
				Error:    "CSV行格式错误: " + err.Error(),
			})
			rowNum++
			continue
		}

		if err := ctl.processCSVRow(record, &result, rowNum); err != nil {
			return nil, err
		}
		rowNum++
	}

	if result.Failed > 0 && result.Success > 0 {
		// 部分失败：使用207 Multi-Status
		return gggin.NewResponseWithStatusCode[any](http.StatusMultiStatus, result), nil
	}

	return gggin.NewResponse[any](result), nil
}

func (ctl *ControllerUser) processCSVRow(record []string, result *BatchRegisterResult, rowNum int) *gggin.HttpError {
	result.Total++

	if len(record) != 6 {
		result.Failed++
		result.Failures = append(result.Failures, BatchFailureInfo{
			Row:      rowNum,
			Username: "",
			Error:    "CSV格式错误：应包含6个字段 (username,surName,givenName,mail,category,role)",
		})
		return nil
	}

	username := strings.TrimSpace(record[0])
	surName := strings.TrimSpace(record[1])
	givenName := strings.TrimSpace(record[2])
	mail := strings.TrimSpace(record[3])
	category := strings.TrimSpace(record[4])
	role := strings.TrimSpace(record[5])

	if username == "" || mail == "" || category == "" || role == "" {
		result.Failed++
		result.Failures = append(result.Failures, BatchFailureInfo{
			Row:      rowNum,
			Username: username,
			Error:    "必填字段不能为空 (username, mail, category, role)",
		})
		return nil
	}

	err := ctl.serviceManager.Register(username, surName, givenName, mail, category, role)
	if err != nil {
		result.Failed++
		result.Failures = append(result.Failures, BatchFailureInfo{
			Row:      rowNum,
			Username: username,
			Error:    err.Error(),
		})
	} else {
		result.Success++
	}

	return nil
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

// ======== 批量操作相关类型定义 ========

// BatchRegisterResult 批量注册结果
type BatchRegisterResult struct {
	Success  int                `json:"success"`
	Failed   int                `json:"failed"`
	Total    int                `json:"total"`
	Failures []BatchFailureInfo `json:"failures"`
}

// BatchFailureInfo 批量操作失败信息
type BatchFailureInfo struct {
	Row      int    `json:"row"`
	Username string `json:"username"`
	Error    string `json:"error"`
}

// RequestBatchModifyCategory 批量修改类别请求体
type RequestBatchModifyCategory struct {
	UserIds  []string `json:"userIds" binding:"required"`
	Category string   `json:"category" binding:"required"`
}

// RequestBatchModifyRole 批量修改角色请求体
type RequestBatchModifyRole struct {
	UserIds []string `json:"userIds" binding:"required"`
	Role    string   `json:"role" binding:"required"`
}

// BatchModifyResult 批量修改结果
type BatchModifyResult struct {
	Success  int                  `json:"success"`
	Failed   int                  `json:"failed"`
	Total    int                  `json:"total"`
	Failures []BatchModifyFailure `json:"failures"`
}

// BatchModifyFailure 批量修改失败信息
type BatchModifyFailure struct {
	UserId string `json:"userId"`
	Error  string `json:"error"`
}

// ======== 批量操作功能实现 ========

// @Summary      批量修改账号类型
// @Description  批量修改指定用户的账号类型。需要 ADMIN 角色权限。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      RequestBatchModifyCategory  true  "批量修改账号类型请求"
// @Success      200  {object} object{data=BatchModifyResult} "全部成功：批量修改结果"
// @Success      207  {object} object{data=BatchModifyResult} "部分失败：批量修改结果"
// @Failure      400  {object} object{data=string} "请求参数错误"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users/category [patch]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleBatchModifyCategory(c *gin.Context) (*gggin.Response[BatchModifyResult], *gggin.HttpError) {
	_, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}

	req, err := gggin.ShouldBindJSON[RequestBatchModifyCategory](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	var result BatchModifyResult
	result.Total = len(req.UserIds)

	for _, uid := range req.UserIds {
		err := ctl.serviceManager.ModifyCategory(uid, req.Category)
		if err != nil {
			result.Failed++
			result.Failures = append(result.Failures, BatchModifyFailure{
				UserId: uid,
				Error:  err.Error(),
			})
		} else {
			result.Success++
		}
	}

	// 根据结果决定返回方式：全部成功返回200，部分失败返回207
	if result.Failed > 0 && result.Success > 0 {
		// 部分失败：使用207 Multi-Status
		return gggin.NewResponseWithStatusCode(http.StatusMultiStatus, result), nil
	}

	// 全部成功或全部失败：使用200
	return gggin.NewResponse(result), nil
}

// @Summary      批量修改角色权限
// @Description  批量修改指定用户的角色权限。需要 ADMIN 角色权限。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      RequestBatchModifyRole  true  "批量修改角色权限请求"
// @Success      200  {object} object{data=BatchModifyResult} "全部成功：批量修改结果"
// @Success      207  {object} object{data=BatchModifyResult} "部分失败：批量修改结果"
// @Failure      400  {object} object{data=string} "请求参数错误"
// @Failure      401  {object} object{data=string} "未授权访问"
// @Failure      403  {object} object{data=string} "权限不足"
// @Failure      500  {object} object{data=string} "服务器内部错误"
// @Router       /users/role [patch]
// @Security     BearerAuth
func (ctl *ControllerUser) HandleBatchModifyRole(c *gin.Context) (*gggin.Response[BatchModifyResult], *gggin.HttpError) {
	_, ok := gggin.Get[*security.GuardResult](c, "guard")
	if !ok {
		return nil, ErrHttpGuardFail
	}

	req, err := gggin.ShouldBindJSON[RequestBatchModifyRole](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	var result BatchModifyResult
	result.Total = len(req.UserIds)

	for _, uid := range req.UserIds {
		err := ctl.serviceManager.GrantRoleByUidAndRoleName(uid, req.Role)
		if err != nil {
			result.Failed++
			result.Failures = append(result.Failures, BatchModifyFailure{
				UserId: uid,
				Error:  err.Error(),
			})
		} else {
			result.Success++
		}
	}

	// 根据结果决定返回方式：全部成功返回200，部分失败返回207
	if result.Failed > 0 && result.Success > 0 {
		// 部分失败：使用207 Multi-Status
		return gggin.NewResponseWithStatusCode(http.StatusMultiStatus, result), nil
	}

	// 全部成功或全部失败：使用200
	return gggin.NewResponse(result), nil
}
