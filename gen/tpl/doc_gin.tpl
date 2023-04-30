package doc

import (
    "encoding/json"
	"github.com/CloverOS/gin-swagger-bootstrap/bootstrap"
	"github.com/CloverOS/goctl-swag/gen"
	"net/http"
    "path/filepath"
    "strings"
)

func WrapHandler() http.HandlerFunc {
	return WrapHandlerWithGroup(Group)
}

func WrapHandlerWithGroupObject(group []gen.SwaggerGroupObject) http.HandlerFunc {
	bytes, err := json.MarshalIndent(group, "", "    ")
	if err != nil {
		return nil
	}
	return WrapHandlerWithGroup(string(bytes))
}

func WrapHandlerWithGroup(group string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = rw.Write([]byte("Method Not Allowed"))
			return
		}
		i := strings.LastIndex(r.URL.Path, "/")
		if i == -1 {
			return
		}
		path := r.URL.Path[i+1:]
		if path == "" {
			path = "index.html"
		}
		switch filepath.Ext(path) {
		case ".html":
			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".css":
			rw.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".js":
			rw.Header().Set("Content-Type", "application/javascript")
		case ".png":
			rw.Header().Set("Content-Type", "image/png")
		case ".json":
			rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
		switch path {
		case "{{.apiJson}}":
			_, _ = rw.Write([]byte(Doc))
		case "group.json":
			_, _ = rw.Write([]byte(group))
		case "index.html":
			readFile, err := bootstrap.ReadFile(path)
			if err != nil {
				bootstrap.Handler.ServeHTTP(rw, r)
				return
			}
			_, _ = rw.Write(readFile)
		default:
			tmp := strings.Replace(r.URL.Path, "/swagger/", "", 1)
			readFile, err := bootstrap.ReadFile(tmp)
			if err != nil {
				bootstrap.Handler.ServeHTTP(rw, r)
				return
			}
			_, _ = rw.Write(readFile)
		}
	}
}

const Doc = `{{.doc}}`

const Group = `{{.group}}`
