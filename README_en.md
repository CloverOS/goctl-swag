# goctl-swag

🌍 *[Simplified Chinese](README.md) ∙ English*  
[Swagger 2.0 Documentation](https://swagger.io/docs/specification/2-0/basic-structure/)

### Introduction

> This project is a plugin based on [go-zero](https://github.com/zeromicro/go-zero), used for generating swagger documentation according to the go-zero's API,
> currently only supporting version 2.0 of swagger. Additionally, the front end page enhancement part is added,
> using [swagger-bootstrap](https://github.com/CloverOS/gin-swagger-bootstrap) for the front end page enhancement,
> inspired by parts of the idea and code from [goctl-swagger](https://github.com/zeromicro/goctl-swagger),
> and using [go-openapi](https://github.com/go-openapi/spec) for parsing and generating swagger.

It also helped me to gain a deeper understanding of go-zero (I have only recently come into contact with go-zero).

## Table of Contents
- [Install Plugin](#install-plugin)
- [Quick Start](#quick-start)
- [Screenshots](#screenshots)
- [goctl-swag Command Usage](#goctl-swag-command-usage)
- [Format Description](#format-description)
    - [General API Information](#general-api-information)
    - [API Declaration](#api-declaration)
    - [API Group Declaration](#api-group-declaration)
    - [Mime Types](#mime-types)
    - [Parameter Types](#parameter-types)
    - [Data Types](#data-types)
- [Usage of Aggregated API Documents](#usage-of-aggregated-api-documents)
- [Acknowledgements](#acknowledgements)

### Install Plugin

```go
go install github.com/CloverOS/goctl-swag/cmd/goctl-swag@latest
```

### Quick Start

***First install [goctl](https://go-zero.dev/docs/goctl/goctl/) and the goctl-swag plugin***

- Take testdata/go-zero/user.api as an example, generate documentation

```shell
git clone https://github.com/CloverOS/goctl-swag
cd goctl-swag
```

- After downloading the project, execute the following command

```shell
goctl api plugin -p goctl-swag="init" -api testdata/user.api -dir testdata/go-zero/user 
```

- This will generate a file `doc_gin.go` in `testdata/go-zero/user/doc`, which is the swagger documentation file containing both the swagger JSON data and the integration of the frontend page enhancement.


- Finally, run the following command to generate a go-zero project (you can refer to [goctl - Plugin Command](https://go-zero.dev/docs/goctl/plugin))

```shell
goctl api go -api testdata/user.api -dir testdata/go-zero/user 
```  

- Adapt for page enhancement
  Modify `test/data/go-zero/user/user.go` as follows, using gin as the router (recommended):

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
    // for gin-swagger-bootstrap doc
    server.AddRoute(rest.Route{
        Method:  http.MethodGet,
        Path:    "/swagger/*any",
        Handler: doc.WrapHandler(),
    })
    fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
    server.Start()
}
```

- Using the default routing

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

- Start the project and access http://localhost:8888/swagger to see the effect

```shell
go mod tidy
go run testdata/go-zero/user/user.go -f testdata/go-zero/user/etc/user.yaml 
```

## Screenshots

![index](/images/swag_index.png)
![settings](/images/swag-settings.png)
![auth](/images/swag_auth.png)
![form](/images/swag_form.png)
![json](/images/swag-json.png)

## goctl-swag Command Usage

```shell
goctl-swag init [command options] [arguments...]
```

| Parameter           | Description                       | Default Value  |
|---------------------|-----------------------------------|----------------|
| --host              | API request address               | localhost:8888 |
| --basepath, --bp    | URL request prefix                | /              |
| --dir, -d           | Folder path for generated doc.go  | ./doc          |
| --webframe, --wb    | Choose which framework to use     | gin            |
| --cover, --cv       | Choose whether to auto overwrite files | true       |
| --router, --rt      | Whether to generate routing information | true      |

## Format Description

### Please refer to [user.api](/testdata/user.api) for examples

### General API Information

#### Example

```
info(
	title: "Example Document"
	desc: "Document Description"
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

| Annotation                 | Explanation                                                   |
|----------------------------|---------------------------------------------------------------|
| title                      | **Required** The name of the application.                     |
| version                    | **Required** The version of the API provided by the application. |
| description                | A short description of the application.                       |
| termsOfService             | The terms of service for the API.                             |
| contact_name               | Contact information for the publicly available API.           |
| contact_url                | URL for the contact information. Must be in URL format.       |
| contact_email              | Email address of the contact person/organization. Must be in email address format. |
| license_name               | **Required** The name of the license used for the API.        |
| license_url                | The URL for the license usedfor the API. Must be in URL format.                                       |
| host                    | The host (hostname or IP) that is running the API.                         |
| BasePath                | The base path on which the API is served.                                  |
| accept                  | A list of MIME types the API can consume. Note that, Accept only affects operations that expect a request body, such as POST, PUT, and PATCH. Values must be as described in "[Mime Types](#mime-types)". |
| produce                 | A list of MIME types the API can produce. Values must be as described in "[Mime Types](#mime-types)". |
| schemes                 | A comma-separated list of the transfer protocols for requests.              |

For strings that cannot use the format like "https://", use the `[$https]` format which will be automatically replaced:

```
[$https] ==> https://
[$http] ==> http://
[$ws] ==> ws://
[$wss] ==> wss://
```

## API Declaration

```api 
@doc(
    summary: "Login"
    description: "Client user login"
)
@handler jsonLoginHandler
post /users/login (JsonRequest) returns (TokenResponse)
```

Or briefly:

```api
@doc "Form Login"
@handler formLoginHandler
post /users/form/login (FormRequest)
```

## API Group Declaration

```api
@server(
    group: login
    open: true
    prefix: /login
    tag: Login Module
)
```
- open: true indicates whether this module is an open interface; by default, it's private (false).

## Request Parameter Declaration

### Json
```api
JsonRequest {
    RequestId string `header:"requestId"`       // request id
    Account   string `json:"account"`           // account
    Password  string `json:"password"`          // password
    SecCode   string `json:"sec_code,optional"` // security code
}
```

### File Upload
```api
UploadRequest {
    MyFile string `form:"myFile,file"` // file
}
```

## Mime Types

`swag` accepts all MIME types with correct format, even if they match `*/*`. Besides these, `swag` also accepts some aliases for certain MIME types, as follows:

| Alias                 | MIME Type                         |
|-----------------------|-----------------------------------|
| json                   | application/json                  |
| xml                    | text/xml                          |
| plain                  | text/plain                        |
| html                   | text/html                         |
| mpfd                   | multipart/form-data               |
| x-www-form-urlencoded  | application/x-www-form-urlencoded |
| json-api               | application/vnd.api+json          |
| json-stream            | application/x-json-stream         |
| octet-stream           | application/octet-stream          |
| png                    | image/png                         |
| jpeg                   | image/jpeg                        |
| gif                    | image/gif                         |

## Parameter Types

- query
- path
- header
- body
- formData

## Data Types

- string (string)
- integer (int, uint, uint32, uint64)
- number (float32)
- boolean (bool)
- user defined struct

## Usage of Aggregated API Documents

> In the generated `doc_gin.go`, there is a method `WrapHandlerWithGroupObject` for amalgamating multiple API documents into one. Here is how to use it:

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
- Here, the URL is filled with different microservices' swag documentation. Afterwards, you can choose the corresponding microservices documentation from the top left corner of the page.

![group](/images/group.png)

### Acknowledgements

[go-zero](https://github.com/zeromicro/go-zero)  
[goctl-swagger](https://github.com/zeromicro/goctl-swagger)  
[go-openapi](https://github.com/go-openapi/spec)