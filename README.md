# goctl-swag

🌍 *简体中文 ∙ [English](README_en.md)*  
[Swagger 2.0 文档](https://swagger.io/docs/specification/2-0/basic-structure/)

### 简介

> 本项目是基于[go-zero](https://github.com/zeromicro/go-zero)写的一个插件，用于根据go-zero的api
> 生成swagger文档，目前只支持swagger2.0版本，另外添加了前端的页面增强部分，
> 使用了[swagger-bootstrap](https://github.com/CloverOS/gin-swagger-bootstrap)作为前端页面部分的增强,
> 借鉴了[goctl-swagger](https://github.com/zeromicro/goctl-swagger)一部分思路和代码，
> 使用了[go-openapi](https://github.com/go-openapi/spec)作为swagger的解析和生成。

也正好借助这个项目对go-zero有了更深的了解（我刚接触go-zero不久）

## 目录
- [安装插件](#安装插件)
- [快速开始](#快速开始)
- [截图](#截图)
- [goctl-swag指令说明](#goctl-swag指令说明)
- [格式说明](#格式说明)
  - [通用API信息](#通用api信息)
  - [API声明](#api声明)
  - [API分组声明](#api分组声明)
  - [Mime类型](#mime类型)
  - [参数类型](#参数类型)
  - [数据类型](#数据类型)
- [聚合Api文档用法](#聚合api文档用法)
- [感谢](#感谢)

### 安装插件

```go
go install github.com/CloverOS/goctl-swag/cmd/goctl-swag@latest
```

### 快速开始

***先安装好[goctl](https://go-zero.dev/docs/goctl/goctl/)以及goctl-swag插件***

- 以testdata/go-zero/user.api为例,进行文档生成

```shell
git clone https://github.com/CloverOS/goctl-swag
cd goctl-swag
```

- 下载本项目后,接着执行以下指令

```shell
goctl api plugin -p goctl-swag="init" -api testdata/user.api -dir testdata/go-zero/user 
```

- 就会在testdata/go-zero/user/doc生成一个doc_gin.go文件，这个文件就是swagger文档的生成文件，里面包含了swagger的json数据和前端页面的增强部分的整合。


- 最后执行以下指令生成go-zero项目（可以参考[goctl-插件指令](https://go-zero.dev/docs/goctl/plugin)）

```shell
goctl api go -api testdata/user.api -dir testdata/go-zero/user 
```  

- 进行页面增强适配
  修改test/data/go-zero/user/user.go如下，使用gin作为router（推荐）

```go
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
    fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
    server.Start()
}
```

- 使用默认路由写法

```go
server.AddRoutes([]rest.Route{
    {
        Method:  http.MethodGet,
        Path:    "/swagger",
        Handler: goctl_swag.Doc("/swagger", "dev"),
    },
    {
        Method: http.MethodGet,
        Path:   "/swagger-json",
        Handler: func (writer http.ResponseWriter, request *http.Request) {
            writer.Header().Set("Content-Type", "application/json; charset=utf-8")
            _, err := writer.Write([]byte(doc.Doc))
            if err != nil {
                httpx.Error(writer, err)
            }
        },
    },
})
```

- 最后启动项目，访问http://localhost:8888/swagger即可看到效果

```shell
go mod tidy
go run testdata/go-zero/user/user.go -f testdata/go-zero/user/etc/user.yaml 
```

## 截图

![index](/images/swag_index.png)
![settings](/images/swag-settings.png)
![auth](/images/swag_auth.png)
![form](/images/swag_form.png)
![json](/images/swag-json.png)

## goctl-swag指令说明

```shell
goctl-swag init [command options] [arguments...]
```

| 参数               | 说明              | 默认值            |
|------------------|-----------------|----------------|
| --host           | API请求地址         | localhost:8888 |
| --basepath, --bp | URL请求前缀         | /              |
| --dir, -d        | 生成的doc.go的文件夹路径 | ./doc          |
| --webframe, --wb | 选择使用哪个框架        | gin            |
| --cover, --cv    | 选择是否自动覆盖文件      | true           |
| --router, --rt   | 是否生成路由信息        | true           |

## 格式说明

### 可以参考[user.api](/testdata/user.api)

### 通用API信息

#### 示例

```
info(
	title: "示例文档"
	desc: "文档描述"
	version: "1.0"
	basePath: "/v1"
	termsOfService: "[$https]swagger.io/terms/"
	contact_name: "khthink"
	contact_url: "[$https]khthink.cn"
	contact_email: "example"
	license_name: "Apache 2.0"
	license_url: "[$https]www.apache.org/licenses/LICENSE-2.0.html"
	BasePath: "/v1"
	accept: "json"
	produce: "json"
	schemes: "http,https"
	securityDefinitions_apikey: "AccessToken"
	in: "header"
	name: "Access-Token"
)
```

| 注释                      | 说明                                                                                            |
|-------------------------|-----------------------------------------------------------------------------------------------|
| title                   | **必填** 应用程序的名称。                                                                               |
| version                 | **必填** 提供应用程序API的版本。                                                                          |
| description             | 应用程序的简短描述。                                                                                    |
| termsOfService          | API的服务条款。                                                                                     |
| contact_name            | 公开的API的联系信息。                                                                                  |
| contact_url             | 联系信息的URL。 必须采用网址格式。                                                                           |
| contact_email           | 联系人/组织的电子邮件地址。 必须采用电子邮件地址的格式。                                                                 |
| license_name            | **必填** 用于API的许可证名称。                                                                           |
| license_url             | 用于API的许可证的URL。 必须采用网址格式。                                                                      |
| host                    | 运行API的主机（主机名或IP地址）。                                                                           |
| BasePath                | 运行API的基本路径。                                                                                   |
| accept                  | API 可以使用的 MIME 类型列表。 请注意，Accept 仅影响具有请求正文的操作，例如 POST、PUT 和 PATCH。 值必须如“[Mime类型](#mime类型)”中所述。 |
| produce                 | API可以生成的MIME类型的列表。值必须如“[Mime类型](#mime类型)”中所述。                                                 |
| schemes                 | 用,分隔的请求的传输协议。                                                                                 |

关于无法使用"https://"类似的这种字符串的问题，可以使用`[$https]`这种方式来解决，会自动替换为`https://`。

```
[$https] ==> https://
[$http] ==> http://
[$ws] ==> ws://
[$wss] ==> wss://
```

## API声明

```api 
@doc(
    summary: "登录"
    description: "客户端用户登录"
)
@handler jsonLoginHandler
post /users/login (JsonRequest) returns (TokenResponse)
```

或者简写

```api
@doc "表单登录"
@handler formLoginHandler
post /users/form/login (FormRequest)
```

## API分组声明

```api
@server(
    group: login
    open: true
    prefix: /login
    tag: 登录模块
)
```
- open: true 代表该模块是否为开放接口，不写的话默认为私有接口（即false）

## 请求参数声明

### Json
```api
JsonRequest {
    RequestId string `header:"requestId"`       //请求id
    Account   string `json:"account"`           //账号
    Password  string `json:"password"`          //密码
    SecCode   string `json:"sec_code,optional"` //安全密码
}
```

### 文件上传
```api
UploadRequest {
    MyFile string `form:"myFile,file"` //文件
}
```

## Mime类型

`swag` 接受所有格式正确的MIME类型, 即使匹配 `*/*`。除此之外，`swag`还接受某些MIME类型的别名，如下所示：

| Alias                 | MIME Type                         |
|-----------------------|-----------------------------------|
| json                  | application/json                  |
| xml                   | text/xml                          |
| plain                 | text/plain                        |
| html                  | text/html                         |
| mpfd                  | multipart/form-data               |
| x-www-form-urlencoded | application/x-www-form-urlencoded |
| json-api              | application/vnd.api+json          |
| json-stream           | application/x-json-stream         |
| octet-stream          | application/octet-stream          |
| png                   | image/png                         |
| jpeg                  | image/jpeg                        |
| gif                   | image/gif                         |


## 参数类型

- query
- path
- header
- body
- formData

## 数据类型

- string (string)
- integer (int, uint, uint32, uint64)
- number (float32)
- boolean (bool)
- user defined struct

## 聚合Api文档用法

> 在生成的doc_gin.go中，有一个WrapHandlerWithGroupObject方法，可以将多个api文档聚合到一个文档中，用法如下：

```go
doc.WrapHandlerWithGroupObject([]gen.SwaggerGroupObject{
    {
        Name:           "user",
        Url:            "../swagger/user.json",
        SwaggerVersion: "2.0",
        Location:       "",
    },
    {
        Name:           "order",
        Url:            "localhost:8889",
        SwaggerVersion: "2.0",
        Location:       "",
    },
})
```
- 这里url填写不同微服务的swag文档数据,之后就可以在页面的左上角中选择对应的微服务文档了
![group](/images/group.png)  


### 感谢

[go-zero](https://github.com/zeromicro/go-zero)  
[goctl-swagger](https://github.com/zeromicro/goctl-swagger)  
[go-openapi](https://github.com/go-openapi/spec)
