package controller

import (
	"errors"
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
	g.POST("", gggin.HandleController(ctl.HandleCreate))
	return ctl
}

type CreateTokenRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary      创建访问令牌
// @Tags         token
// @Accept       json
// @Produce      json
// @Param        body  body      CreateTokenRequest  true  "创建令牌请求"
// @Success      200  {object} object{data=string} "返回访问令牌"
// @Router       /api/tokens [post]
func (ctl *ControllerToken) HandleCreate(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	req, err := gggin.ShouldBindJSON[CreateTokenRequest](c)
	if err != nil {
		return nil, gggin.NewHttpError(http.StatusBadRequest, err.Error())
	}

	token, err := ctl.serviceManager.Authenticate(req.Username, req.Password)
	if err != nil {
		var httpErr *gggin.HttpError
		if errors.As(err, &httpErr) {
			return nil, httpErr
		}
		return nil, gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	return gggin.NewResponse(token), nil
}
