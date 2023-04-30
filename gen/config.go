package gen

import (
	"encoding/json"
	"fmt"
)

// Config presents Gen configurations.
type Config struct {
	//Host api request address
	Host string

	//Dir the folder path of the generated doc.go
	Dir string

	//BasePath url request prefix
	BasePath string

	//WebFrame such as http,gin
	WebFrame string

	//Cover auto cover generation
	Cover bool
}

// FileGenConfig file creation configuration
type FileGenConfig struct {
	dir             string
	subDir          string
	filename        string
	templateName    string
	category        string
	templateFile    string
	builtinTemplate string
	data            any
}

func (c *Config) fmtString() {
	marshal, _ := json.Marshal(c)
	fmt.Println(string(marshal))
}

func (c *FileGenConfig) fmtString() {
	marshal, _ := json.Marshal(c)
	fmt.Println(string(marshal))
}
