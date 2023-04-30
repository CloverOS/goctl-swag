package login

import (
	"context"

	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/svc"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JsonLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJsonLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JsonLoginLogic {
	return &JsonLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JsonLoginLogic) JsonLogin(req *types.JsonRequest) (resp *types.TokenResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
