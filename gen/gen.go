package gen

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CloverOS/goctl-swag/util"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"go/format"
	"io"
	"io/fs"
	"log"
	"os"
	"text/template"
)

const (
	Gin = "gin"
	Chi = "chi"
	Mux = "mux"
)

//go:embed tpl/*
var templateFiles embed.FS

// Gen presents a generate tool for swag.
type Gen struct {
	config     *Config
	json       func(data interface{}) ([]byte, error)
	jsonIndent func(data interface{}) ([]byte, error)
	debug      Debugger
}

// Debugger is the interface that wraps the basic Printf method.
type Debugger interface {
	Printf(format string, v ...interface{})
}

func New(config *Config) *Gen {
	return &Gen{
		config: config,
		json:   json.Marshal,
		jsonIndent: func(data interface{}) ([]byte, error) {
			return json.MarshalIndent(data, "", "    ")
		},
		debug: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (g *Gen) DoGen() error {
	p, err := plugin.NewPlugin()
	if err != nil {
		return err
	}
	fmt.Println("start generating documentation ...")
	swag, err := NewParser(p).parseSwag(g)
	if err != nil {
		return err
	}
	indent, err := g.jsonIndent(&swag)
	if err != nil {
		return err
	}
	var filename string
	var templateFile string
	switch g.config.WebFrame {
	case Gin:
		filename = "doc_gin.go"
		templateFile = "tpl/doc_gin.tpl"
	//case Chi:
	//	filename = "doc_chi.go"
	//	templateFile = "tpl/doc_chi.tpl"
	//case Mux:
	//	filename = "doc_mux.go"
	//	templateFile = "tpl/doc_mux.tpl"
	default:
		return errors.New("unknown web frame")
	}
	return g.genFile(FileGenConfig{
		dir:          p.Dir,
		subDir:       g.config.Dir,
		filename:     filename,
		templateName: "docTemplate",
		category:     "doc",
		templateFile: templateFile,
		data: map[string]string{
			"doc":     string(indent),
			"group":   g.getGroup(p),
			"apiJson": p.Api.Service.Name + ".json",
		},
	})
}

func (g *Gen) genFile(fgc FileGenConfig) error {
	fp, created, err := util.MaybeCreateFile(fgc.dir, fgc.subDir, fgc.filename, g.config.Cover)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer func(fp *os.File) {
		_ = fp.Close()
	}(fp)

	var text string
	if len(fgc.category) == 0 || len(fgc.templateFile) == 0 {
		text = fgc.builtinTemplate
	} else {
		// Read template files from the embedded filesystem
		file, err := templateFiles.Open(fgc.templateFile)
		if err != nil {
			return err
		}
		defer func(file fs.File) {
			_ = file.Close()
		}(file)
		// Read template file content
		content, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		text = string(content)
	}

	t := template.Must(template.New(fgc.templateName).Parse(text))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, fgc.data)
	if err != nil {
		return err
	}
	b := buffer.String()
	var code string
	ret, err := format.Source([]byte(b))
	if err != nil {
		code = b
	} else {
		code = string(ret)
	}
	_, err = fp.WriteString(code)
	return err
}

func (g *Gen) getGroup(p *plugin.Plugin) string {
	indent, err := g.jsonIndent([]SwaggerGroupObject{
		{
			Name:           p.Api.Service.Name,
			Url:            "../swagger/" + p.Api.Service.Name + ".json",
			SwaggerVersion: "2.0",
			Location:       "",
		},
	})
	if err != nil {
		return ""
	}
	return string(indent)
}
