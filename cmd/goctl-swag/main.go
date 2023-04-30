package main

import (
	"fmt"
	"github.com/CloverOS/goctl-swag/gen"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
)

const Version = "v0.0.8"

const (
	HostFlag     = "host"
	DirFlag      = "dir"
	BasePathFlag = "basepath"
	WebFrameFlag = "webframe"
	CoverFlag    = "cover"
)

var (
	initFlags = []cli.Flag{
		&cli.StringFlag{
			Name:  HostFlag,
			Value: "localhost:8888",
			Usage: "api request address",
		},
		&cli.StringFlag{
			Name:    BasePathFlag,
			Aliases: []string{"bp"},
			Value:   "/",
			Usage:   "url request prefix",
		},
		&cli.StringFlag{
			Name:    DirFlag,
			Aliases: []string{"d"},
			Value:   "./doc",
			Usage:   "the folder path of the generated doc.go",
		},
		&cli.StringFlag{
			Name:    WebFrameFlag,
			Aliases: []string{"wb"},
			Value:   "gin",
			Usage:   "select which webframe used",
		},
		&cli.BoolFlag{
			Name:    CoverFlag,
			Aliases: []string{"cv"},
			Value:   true,
			Usage:   "choose whether to automatically overwrite files",
		},
	}
)

func initAction(ctx *cli.Context) error {
	return gen.New(&gen.Config{
		Host:     ctx.String(HostFlag),
		Dir:      ctx.String(DirFlag),
		BasePath: ctx.String(BasePathFlag),
		WebFrame: ctx.String(WebFrameFlag),
		Cover:    ctx.Bool(CoverFlag),
	}).DoGen()
}

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s %s/%s", Version, runtime.GOOS, runtime.GOARCH)
	app.Usage = "a plugin of goctl to generate swagger document with extra webui"
	app.Commands = []*cli.Command{
		{
			Name:   "init",
			Usage:  "generates swagger document",
			Action: initAction,
			Flags:  initFlags,
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("goctl-swagger: %+v\n", err)
	}
}
