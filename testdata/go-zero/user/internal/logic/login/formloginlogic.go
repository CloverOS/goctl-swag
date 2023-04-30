package login

import (
	"context"

	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/svc"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FormLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFormLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FormLoginLogic {
	return &FormLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FormLoginLogic) FormLogin(req *types.FormRequest) error {
	// todo: add your logic here and delete this line

	return nil
}
