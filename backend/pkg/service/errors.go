package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dsx137/gg-gin/pkg/gggin"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidRole     = errors.New("invalid role")
	ErrInvalidOu       = errors.New("invalid organizational unit")
	ErrInvalidCreds    = errors.New("invalid credentials")
	ErrConflict        = errors.New("conflict")
	ErrIllegalPassword = errors.New("illegal password")
	ErrWeakPassword    = errors.New("weak password")
)

type ServiceError struct {
	Err     error
	Message string
}

func NewError(err error, message string) *ServiceError {
	return &ServiceError{Err: err, Message: message}
}

func (e *ServiceError) Error() string { return e.Message }

func MapErrorToHttp(err error) *gggin.HttpError {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return gggin.NewHttpError(http.StatusNotFound, fmt.Sprintf("用户不存在: %s", err.Error()))
	case errors.Is(err, ErrUserExists):
		return gggin.NewHttpError(http.StatusConflict, fmt.Sprintf("用户已存在: %s", err.Error()))
	case errors.Is(err, ErrInvalidEmail):
		return gggin.NewHttpError(http.StatusBadRequest, fmt.Sprintf("邮箱格式不合法: %s", err.Error()))
	case errors.Is(err, ErrInvalidOu):
		return gggin.NewHttpError(http.StatusBadRequest, fmt.Sprintf("账号类型不合法: %s", err.Error()))
	case errors.Is(err, ErrInvalidRole):
		return gggin.NewHttpError(http.StatusBadRequest, fmt.Sprintf("角色不合法: %s", err.Error()))
	case errors.Is(err, ErrInvalidCreds):
		return gggin.NewHttpError(http.StatusUnauthorized, fmt.Sprintf("无效的凭证: %s", err.Error()))
	case errors.Is(err, ErrIllegalPassword):
		return gggin.NewHttpError(http.StatusBadRequest, fmt.Sprintf("密码不符合要求: %s", err.Error()))
	case errors.Is(err, ErrWeakPassword):
		return gggin.NewHttpError(http.StatusBadRequest, fmt.Sprintf("密码强度不够: %s", err.Error()))
	default:
		return gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}
}
