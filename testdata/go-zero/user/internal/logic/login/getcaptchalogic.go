package login

import (
	"context"

	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/svc"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetcaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetcaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetcaptchaLogic {
	return &GetcaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetcaptchaLogic) Getcaptcha(req *types.GetCaptchaRequest) error {
	// todo: add your logic here and delete this line

	return nil
}
