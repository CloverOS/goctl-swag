package main

import (
	"flag"
	"fmt"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/doc"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/config"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/handler"
	"github.com/CloverOS/goctl-swag/testdata/go-zero/user/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/zero-contrib/router/gin"
	"net/http"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	r := gin.NewRouter()
	server := rest.MustNewServer(c.RestConf, rest.WithRouter(r))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	//for gin-swagger-bootstrap doc
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/swagger/*any",
		Handler: doc.WrapHandler(),
	})

	////for simple swagger doc
	//server.AddRoutes([]rest.Route{
	//	{
	//		Method:  http.MethodGet,
	//		Path:    "/swagger",
	//		Handler: goctl_swag.Doc("/swagger", "dev"),
	//	},
	//	{
	//		Method: http.MethodGet,
	//		Path:   "/swagger-json",
	//		Handler: func(writer http.ResponseWriter, request *http.Request) {
	//			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	//			_, err := writer.Write([]byte(doc.Doc))
	//			if err != nil {
	//				httpx.Error(writer, err)
	//			}
	//		},
	//	},
	//})
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
