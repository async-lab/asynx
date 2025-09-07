package controller

import (
	"net/http"

	"asynclab.club/asynx/backend/pkg/service"
	"github.com/dsx137/gg-gin/pkg/gggin"
	"github.com/gin-gonic/gin"
)

type ControllerToken struct {
	serviceManager *service.ServiceManager
}

func NewControllerTokens(g *gin.RouterGroup, serviceManager *service.ServiceManager) *ControllerToken {
	ctl := &ControllerToken{serviceManager: serviceManager}
	g.POST("", gggin.ToGinHandler(ctl.HandleCreate))
	return ctl
}

type CreateTokenRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary      创建访问令牌
// @Description  通过用户名和密码验证用户身份并生成访问令牌
// @Tags         tokens
// @Accept       json
// @Produce      json
// @Param        body  body      CreateTokenRequest  true  "创建令牌请求"
// @Success      200   {object}  object{data=string} "返回访问令牌"
// @Failure      400   {object}  object{data=string} "请求参数错误"
// @Failure      401   {object}  object{data=string} "用户名或密码错误"
// @Failure      500   {object}  object{data=string} "服务器内部错误"
// @Router       /tokens [post]
func (ctl *ControllerToken) HandleCreate(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	req, err := gggin.ShouldBindJSON[CreateTokenRequest](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	token, err := ctl.serviceManager.Authenticate(req.Username, req.Password)
	if err != nil {
		return nil, service.MapErrorToHttp(err)
	}
	return gggin.NewResponse(token), nil
}
