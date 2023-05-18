package doc

import (
    "github.com/CloverOS/goctl-swag/doc"
)

var Route = RouteImpl{}

type RouteImpl struct {
}

func (r RouteImpl) GetRouteInfos() []doc.RouteInfos {
    return []doc.RouteInfos{ {{range .}}
    {
        BasePath:   "{{.BasePath}}",
        HandlerFun: "{{.HandlerFun}}",
        Method:     "{{.Method}}",
        Path:       "{{.Path}}",
        Public:     {{.Public}},
        RouteGroup: doc.RouteGroup{GroupName: "{{.RouteGroup.GroupName}}"},
        Summary:    "{{.Summary}}",
    },
    {{end}}
    }
}
