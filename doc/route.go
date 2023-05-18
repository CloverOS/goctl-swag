package doc

type RouteGroup struct {
	GroupName string `json:"name"`
}

type RouteInfos struct {
	Method     string `json:"method"`      //  method
	Path       string `json:"path"`        //  path
	BasePath   string `json:"base_path"`   //  BasePath
	HandlerFun string `json:"handler_fun"` //  handlerFun
	Summary    string `json:"summary"`     //  Summary
	Public     bool   `json:"public"`      //  is public router
	RouteGroup
}

type RouteGen interface {
	GetRouteInfos() []RouteInfos
}
