package login

import (
	"net/http"

	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/logic/login"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/svc"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FormLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FormRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := login.NewFormLoginLogic(r.Context(), svcCtx)
		err := l.FormLogin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJson(w, types.Result{Code: 200, Msg: "Success", Data: nil})
		}
	}
}
