package controller

import (
	"net/http"

	"github.com/dsx137/gg-gin/pkg/gggin"
)

var (
	ErrHttpForceForbidden = gggin.NewHttpError(http.StatusForbidden, "WHAT ARE YOU DOING?")
	ErrHttpGuardFail      = gggin.NewHttpError(http.StatusInternalServerError, "获取用户信息失败")
)
