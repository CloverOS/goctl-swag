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
		case "user.json":
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

const Doc = `{
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
    "schemes": [
        "http,https"
    ],
    "swagger": "2.0",
    "info": {
        "description": "用户服务文档",
        "title": "用户服务接口文档",
        "termsOfService": "https://swagger.io/terms/",
        "contact": {
            "name": "khthink",
            "url": "https://khthink.cn",
            "email": "example"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "https://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "localhost:8888",
    "basePath": "/",
    "paths": {
        "/fileserver/file/download": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文件管理"
                ],
                "summary": "文件下载",
                "operationId": "downloadHandler",
                "parameters": [
                    {
                        "type": "string",
                        "description": "文件夹",
                        "name": "fileDir",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "文件名",
                        "name": "fileName",
                        "in": "query",
                        "required": true
                    }
                ]
            }
        },
        "/fileserver/file/upload": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文件管理"
                ],
                "summary": "文件上传",
                "operationId": "upload",
                "parameters": [
                    {
                        "type": "file",
                        "description": "文件",
                        "name": "myFile",
                        "in": "formData",
                        "required": true
                    }
                ]
            }
        },
        "/login/get/captcha/:key": {
            "get": {
                "security": [
                    {
                        "AccessToken": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "登录模块"
                ],
                "summary": "获取验证码",
                "operationId": "getcaptchaHandler",
                "parameters": [
                    {
                        "type": "string",
                        "description": "验证码key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    }
                ]
            }
        },
        "/login/users/form/login": {
            "post": {
                "security": [
                    {
                        "AccessToken": []
                    }
                ],
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "登录模块"
                ],
                "summary": "表单登录",
                "operationId": "formLoginHandler",
                "parameters": [
                    {
                        "type": "string",
                        "description": "账号",
                        "name": "account",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "密码",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "安全密码",
                        "name": "sec_code",
                        "in": "formData"
                    }
                ]
            }
        },
        "/login/users/login": {
            "post": {
                "security": [
                    {
                        "AccessToken": []
                    }
                ],
                "description": "客户端用户登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "登录模块"
                ],
                "summary": "登录",
                "operationId": "jsonLoginHandler",
                "parameters": [
                    {
                        "type": "string",
                        "description": "请求id",
                        "name": "requestId",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "  ",
                        "name": "JsonRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/JsonRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/TokenResponse"
                        }
                    }
                }
            }
        },
        "/user/edit/myheadimg": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "编辑自己的头像组",
                "operationId": "editMyHeadImgs",
                "parameters": [
                    {
                        "description": "  ",
                        "name": "EditUserHeadImg",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/EditUserHeadImg"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/MyResponse"
                        }
                    }
                }
            }
        },
        "/user/list/:page/:count": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "获取用户列表",
                "operationId": "userListHandler",
                "parameters": [
                    {
                        "type": "int",
                        "description": "页数:1开始",
                        "name": "page",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "int",
                        "description": "每页数量",
                        "name": "count",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "  ",
                        "name": "PageRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/PageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/UserListResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "EditUserHeadImg": {
            "type": "object",
            "title": "EditUserHeadImg",
            "properties": {
                "headImg": {
                    "description": "头像图片组",
                    "type": "array",
                    "required": [
                        "true"
                    ],
                    "items": {
                        "type": "string"
                    }
                },
                "userId": {
                    "description": "用户id",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                }
            }
        },
        "JsonRequest": {
            "type": "object",
            "title": "JsonRequest",
            "properties": {
                "account": {
                    "description": "账号",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                },
                "password": {
                    "description": "密码",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                },
                "requestId": {
                    "description": "请求id",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                },
                "sec_code": {
                    "description": "安全密码",
                    "type": "string"
                }
            }
        },
        "MyResponse": {
            "type": "object",
            "title": "MyResponse",
            "properties": {
                "code": {
                    "description": "错误码",
                    "type": "int",
                    "required": [
                        "true"
                    ]
                },
                "msg": {
                    "description": "错误消息",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                }
            }
        },
        "PageRequest": {
            "type": "object",
            "title": "PageRequest",
            "properties": {
                "nickName": {
                    "description": "昵称",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                }
            }
        },
        "TestNest": {
            "type": "object",
            "title": "TestNest",
            "properties": {
                "test": {
                    "description": "测试",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                }
            }
        },
        "TokenResponse": {
            "type": "object",
            "title": "TokenResponse",
            "properties": {
                "token": {
                    "description": "token",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                }
            }
        },
        "User": {
            "type": "object",
            "title": "User",
            "properties": {
                "account": {
                    "description": "账号",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                },
                "headImg": {
                    "description": "头像图片组",
                    "type": "array",
                    "required": [
                        "true"
                    ],
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "description": "姓名",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                },
                "password": {
                    "description": "密码",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                }
            }
        },
        "UserListResponse": {
            "type": "object",
            "title": "UserListResponse",
            "properties": {
                "userList": {
                    "description": "用户列表",
                    "type": "array",
                    "required": [
                        "true"
                    ],
                    "items": {
                        "type": "object",
                        "$ref": "#/definitions/User"
                    }
                }
            }
        },
        "UserVip": {
            "type": "object",
            "title": "UserVip",
            "properties": {
                "vipId": {
                    "description": "vipid",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                },
                "vipLevel": {
                    "description": "vip等级",
                    "type": "string",
                    "required": [
                        "true"
                    ]
                }
            }
        }
    },
    "securityDefinitions": {
        "apiKey": {
            "description": "AccessToken",
            "type": "apiKey",
            "name": "AccessToken",
            "in": "header"
        }
    }
}`

const Group = `[
    {
        "name": "user",
        "url": "../swagger/user.json",
        "swaggerVersion": "2.0",
        "location": ""
    }
]`
