package gen

import (
	"reflect"
)

var swaggerMapTypes = map[string]reflect.Kind{
	"string":   reflect.String,
	"*string":  reflect.String,
	"int":      reflect.Int,
	"*int":     reflect.Int,
	"uint":     reflect.Uint,
	"*uint":    reflect.Uint,
	"int8":     reflect.Int8,
	"*int8":    reflect.Int8,
	"uint8":    reflect.Uint8,
	"*uint8":   reflect.Uint8,
	"int16":    reflect.Int16,
	"*int16":   reflect.Int16,
	"uint16":   reflect.Uint16,
	"*uint16":  reflect.Uint16,
	"int32":    reflect.Int,
	"*int32":   reflect.Int,
	"uint32":   reflect.Int,
	"*uint32":  reflect.Int,
	"uint64":   reflect.Int64,
	"*uint64":  reflect.Int64,
	"int64":    reflect.Int64,
	"*int64":   reflect.Int64,
	"[]string": reflect.Slice,
	"[]int":    reflect.Slice,
	"[]int64":  reflect.Slice,
	"[]int32":  reflect.Slice,
	"[]uint32": reflect.Slice,
	"[]uint64": reflect.Slice,
	"bool":     reflect.Bool,
	"*bool":    reflect.Bool,
	"struct":   reflect.Struct,
	"*struct":  reflect.Struct,
	"float32":  reflect.Float32,
	"*float32": reflect.Float32,
	"float64":  reflect.Float64,
	"*float64": reflect.Float64,
}

var MineTypeMap = map[string]string{
	"json":        "application/json",
	"xml":         "application/xml",
	"html":        "text/html",
	"text":        "text/plain",
	"form":        "application/x-www-form-urlencoded",
	"mpfd":        "multipart/form-data",
	"jsonapi":     "application/vnd.api+json",
	"jsonstream":  "application/x-json-stream",
	"octetstream": "application/octet-stream",
	"png":         "image/png",
	"jpeg":        "image/jpeg",
	"gif":         "image/gif",
}

var InTypeMap = map[string]string{
	"json":     "body",
	"path":     "path",
	"postform": "formData",
	"mpfd":     "formData",
	"header":   "header",
	"getform":  "query",
}

// SwaggerGroupObject Microservices Documentation Group
type SwaggerGroupObject struct {
	Name           string `json:"name"`
	Url            string `json:"url"`
	SwaggerVersion string `json:"swaggerVersion"`
	Location       string `json:"location"`
}

type RouteGroup struct {
	GroupName string
}

type RouteInfos struct {
	Method     string //  method
	Path       string //  path
	BasePath   string //  BasePath
	HandlerFun string //  handlerFun
	Summary    string //  Summary
	Public     bool   //  is public router
	RouteGroup
}
