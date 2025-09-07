package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dsx137/gg-gin/pkg/gggin"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrExists       = errors.New("already exists")
	ErrInvalid      = errors.New("invalid objet")
	ErrWeakPassword = errors.New("weak password")
)

type ServiceError struct {
	Err     error
	Message string
}

func WrapError(err error, message string) *ServiceError {
	return &ServiceError{Err: err, Message: message}
}

func (e *ServiceError) Error() string { return e.Message }

func (e *ServiceError) Unwrap() error { return e.Err }

func MapErrorToHttp(err error) *gggin.HttpError {
	switch {
	case errors.Is(err, ErrNotFound):
		return gggin.NewHttpError(http.StatusNotFound, fmt.Sprintf("对象不存在: %s", err.Error()))
	case errors.Is(err, ErrExists):
		return gggin.NewHttpError(http.StatusConflict, fmt.Sprintf("对象已存在: %s", err.Error()))
	case errors.Is(err, ErrInvalid):
		return gggin.NewHttpError(http.StatusBadRequest, fmt.Sprintf("无效的对象: %s", err.Error()))
	case errors.Is(err, ErrWeakPassword):
		return gggin.NewHttpError(http.StatusBadRequest, fmt.Sprintf("密码强度不够: %s", err.Error()))
	default:
		return gggin.NewHttpError(http.StatusInternalServerError, err.Error())
	}
}
