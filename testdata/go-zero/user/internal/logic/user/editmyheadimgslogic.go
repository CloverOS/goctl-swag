package user

import (
	"context"

	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/svc"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditMyHeadImgsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditMyHeadImgsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditMyHeadImgsLogic {
	return &EditMyHeadImgsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditMyHeadImgsLogic) EditMyHeadImgs(req *types.EditUserHeadImg) (resp *types.MyResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
