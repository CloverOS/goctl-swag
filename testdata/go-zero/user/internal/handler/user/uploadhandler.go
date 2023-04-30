package user

import (
	"net/http"

	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/logic/user"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/svc"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUploadLogic(r.Context(), svcCtx)
		err := l.Upload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJson(w, types.Result{Code: 200, Msg: "Success", Data: nil})
		}
	}
}
