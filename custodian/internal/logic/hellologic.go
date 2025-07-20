package logic

import (
	"context"
	"fmt"

	"asynclab.club/asynx/custodian/internal/svc"
	"asynclab.club/asynx/custodian/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type HelloLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloLogic {
	return &HelloLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HelloLogic) Hello(req *types.Request) (resp *types.Response, err error) {
	resp = &types.Response{
		Message: fmt.Sprintf("Hello, %s", req.Name),
	}

	return
}
