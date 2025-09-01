package controller

import (
	"github.com/dsx137/gg-gin/pkg/gggin"
	"github.com/gin-gonic/gin"
)

type ControllerHello struct{}

func NewControllerHello(g *gin.RouterGroup) *ControllerHello {
	ctl := &ControllerHello{}
	g.GET("", gggin.ToGinHandler(ctl.HandleHello))
	return ctl
}

// @Summary      打招呼
// @Tags         index
// @Accept       json
// @Produce      json
// @Success      200  {object} object{data=string} "Hello, AsyncLab"
// @Router       /hello [get]
func (ctl *ControllerHello) HandleHello(c *gin.Context) (*gggin.Response[string], *gggin.HttpError) {
	return gggin.NewResponse("Hello, AsyncLab!"), nil
}
