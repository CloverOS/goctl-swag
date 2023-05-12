package gen

import (
	"errors"
	openapi "github.com/go-openapi/spec"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"log"
	"net/http"
	"os"
	"strings"
)

type Parser struct {
	Swagger   *openapi.Swagger
	Plugin    *plugin.Plugin
	RouterMap map[string]RouteInfos
	debug     Debugger
}

func NewParser(p *plugin.Plugin) *Parser {
	return &Parser{
		Plugin:    p,
		RouterMap: make(map[string]RouteInfos),
		debug:     log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (ps *Parser) parseSwag(g *Gen) (openapi.Swagger, error) {
	info := ps.getSwaggerInfoObject(ps.Plugin)

	basePath := g.config.BasePath
	// info() priority is higher than cmd parameter
	if ps.Plugin.Api.Info.Properties["basePath"] != "" {
		basePath = ps.Plugin.Api.Info.Properties["basePath"]
	}
	ps.Swagger = &openapi.Swagger{
		SwaggerProps: openapi.SwaggerProps{
			Schemes:  openapi.StringOrArray{ps.Plugin.Api.Info.Properties["schemes"]},
			Consumes: openapi.StringOrArray{MineTypeMap[ps.Plugin.Api.Info.Properties["accept"]]},
			Produces: openapi.StringOrArray{MineTypeMap[ps.Plugin.Api.Info.Properties["produce"]]},
			Swagger:  "2.0",
			Info:     &info,
			Host:     g.config.Host,
			BasePath: basePath,
			Paths: &openapi.Paths{
				Paths: make(map[string]openapi.PathItem),
			},
			Definitions:         make(openapi.Definitions),
			Parameters:          make(map[string]openapi.Parameter),
			SecurityDefinitions: make(openapi.SecurityDefinitions),
		},
	}

	ps.getSecurityDefinitions()

	err := ps.getPaths()
	if err != nil {
		return openapi.Swagger{}, err
	}

	ps.getDefinitions()

	return *ps.Swagger, nil
}

func (ps *Parser) getPaths() error {
	groups := ps.Plugin.Api.Service.Groups
	for _, group := range groups {
		open := false
		if group.GetAnnotation("open") == "true" {
			open = true
		}
		tag := group.GetAnnotation("tag")
		for _, route := range group.Routes {
			pathItem := openapi.PathItem{}
			routeInfo := RouteInfos{}
			err := ps.getOperation(&pathItem, &route, open, tag, &routeInfo)
			if err != nil {
				return err
			}
			path := group.GetAnnotation("prefix") + route.Path
			if path[0] != '/' {
				path = "/" + path
			}
			ps.Swagger.Paths.Paths[path] = pathItem
			routeInfo.BasePath = ps.Swagger.BasePath
			routeInfo.Path = path
			routeInfo.Public = open
			routeInfo.RouteGroup = RouteGroup{
				GroupName: group.GetAnnotation("tag"),
			}
			ps.RouterMap[path] = routeInfo
		}
	}
	return nil
}

func (ps *Parser) getDefinitions() {
	for _, tp := range ps.Plugin.Api.Types {
		defineStruct, ok := tp.(spec.DefineStruct)
		if !ok {
			continue
		}

		if len(defineStruct.Members) < 0 {
			continue
		}
		if defineStruct.Members[0].IsBodyMember() {
			ps.Swagger.Definitions[defineStruct.RawName] = ps.getDefinition(defineStruct)
		}
	}
}

func (ps *Parser) getDefinition(defineStruct spec.DefineStruct) openapi.Schema {
	//Model
	schema := openapi.Schema{}
	schema.Type = openapi.StringOrArray{"object"}
	schema.Properties = make(map[string]openapi.Schema)
	schema.Title = defineStruct.Name()
	if len(defineStruct.Comments()) > 0 {
		schema.Description = defineStruct.Comments()[0]
	}

	ps.debug.Printf("getDefinition:%s", defineStruct.Name())

	for _, member := range defineStruct.GetNonBodyMembers() {
		if tempDefineStruct, ok := member.Type.(spec.DefineStruct); ok {
			inlineSchema := ps.getDefinition(tempDefineStruct)
			for k, v := range inlineSchema.Properties {
				schema.Properties[k] = v
			}
		}
	}

	for _, member := range defineStruct.Members {
		//Model's property schema
		propSchema := openapi.Schema{}
		key := ""
		for _, tag := range member.Tags() {
			if tag.Key == "header" || tag.Key == "form" || tag.Key == "json" {
				key = tag.Name
				break
			}
		}
		//if key doesn't exist, skip
		if key == "" {
			continue
		}
		propSchema.Description = strings.ReplaceAll(member.GetComment(), "//", "")
		ps.getParamType(&propSchema, member)
		if !member.IsOptional() {
			propSchema.Required = openapi.StringOrArray{"true"}
		}
		schema.Properties[key] = propSchema
	}

	return schema
}

func (ps *Parser) getSecurityDefinitions() {
	apikey, ok := ps.Plugin.Api.Info.Properties["securityDefinitions_apikey"]
	if ok {
		newSecDefValue := openapi.SecurityScheme{}
		newSecDefValue.Name = apikey
		newSecDefValue.Description = apikey
		newSecDefValue.Type = "apiKey"
		newSecDefValue.In = ps.Plugin.Api.Info.Properties["in"]
		ps.Swagger.SecurityDefinitions["apiKey"] = &newSecDefValue
	}
	basic, ok := ps.Plugin.Api.Info.Properties["securityDefinitions_basic"]
	if ok {
		newSecDefValue := openapi.SecurityScheme{}
		newSecDefValue.Description = basic
		newSecDefValue.Type = "basic"
		ps.Swagger.SecurityDefinitions["basic"] = &newSecDefValue
	}
	//todo oauth2
	oauth2Application, ok := ps.Plugin.Api.Info.Properties["securityDefinitions_oauth2_application"]
	if ok {
		newSecDefValue := openapi.SecurityScheme{}
		newSecDefValue.Description = oauth2Application
		newSecDefValue.Type = "oauth2.application"
		newSecDefValue.TokenURL = ps.Plugin.Api.Info.Properties["tokenUrl"]
		newSecDefValue.Scopes = make(map[string]string)
		ps.Swagger.SecurityDefinitions["oauth2.application"] = &newSecDefValue
	}
}

func (ps *Parser) getSwaggerInfoObject(p *plugin.Plugin) openapi.Info {
	p.Api.Info.Properties = ps.trimMap(p.Api.Info.Properties)
	return openapi.Info{
		InfoProps: openapi.InfoProps{
			Description:    p.Api.Info.Properties["desc"],
			Title:          p.Api.Info.Properties["title"],
			TermsOfService: p.Api.Info.Properties["termsOfService"],
			Contact: &openapi.ContactInfo{
				ContactInfoProps: openapi.ContactInfoProps{
					Name:  p.Api.Info.Properties["contact_name"],
					URL:   p.Api.Info.Properties["contact_url"],
					Email: p.Api.Info.Properties["contact_email"],
				},
			},
			License: &openapi.License{
				LicenseProps: openapi.LicenseProps{
					Name: p.Api.Info.Properties["license_name"],
					URL:  p.Api.Info.Properties["license_url"],
				},
			},
			Version: p.Api.Info.Properties["version"],
		},
	}
}

func (ps *Parser) trim(data string) string {
	return strings.Trim(data, "\"")
}

func (ps *Parser) trimMap(data map[string]string) map[string]string {
	for k, v := range data {
		t := ps.trim(v)
		data[k] = ps.transferSpecChar(t)
	}
	return data
}

func (ps *Parser) transferSpecChar(data string) string {
	data = strings.ReplaceAll(data, "[$https]", "https://")
	data = strings.ReplaceAll(data, "[$http]", "http://")
	data = strings.ReplaceAll(data, "[$ws]", "ws://")
	data = strings.ReplaceAll(data, "[$wss]", "wss://")
	return data
}

func (ps *Parser) getOperation(pathItem *openapi.PathItem, route *spec.Route,
	open bool, tag string, routeInfo *RouteInfos) error {
	summary := ""
	desc := ""
	if route.AtDoc.Text == "" {
		summary = ps.trim(route.AtDoc.Properties["summary"])
		desc = ps.trim(route.AtDoc.Properties["description"])
	} else {
		summary = ps.trim(route.AtDoc.Text)
	}
	operation := openapi.NewOperation(route.Handler)
	operation.Summary = summary
	operation.Description = desc
	operation.Tags = []string{ps.trim(tag)}

	routeInfo.Method = route.Method
	routeInfo.HandlerFun = route.Handler
	routeInfo.Summary = operation.Summary

	switch strings.ToUpper(route.Method) {
	case http.MethodGet:
		pathItem.Get = operation
		ps.genOperation(pathItem.Get, route, open)
	case http.MethodPost:
		pathItem.Post = operation
		ps.genOperation(pathItem.Post, route, open)
	case http.MethodDelete:
		pathItem.Delete = operation
		ps.genOperation(pathItem.Delete, route, open)
	case http.MethodHead:
		pathItem.Head = operation
		ps.genOperation(pathItem.Head, route, open)
	case http.MethodPut:
		pathItem.Put = operation
		ps.genOperation(pathItem.Put, route, open)
	case http.MethodOptions:
		pathItem.Options = operation
		ps.genOperation(pathItem.Options, route, open)
	case http.MethodPatch:
		pathItem.Patch = operation
		ps.genOperation(pathItem.Patch, route, open)
	default:
		return errors.New("route method not support")
	}
	return nil
}

func (ps *Parser) genOperation(op *openapi.Operation, route *spec.Route, open bool) {
	if !open {
		for _, s := range ps.Swagger.SecurityDefinitions {
			temp := make(map[string][]string)
			temp[s.Name] = []string{}
			op.Security = append(op.Security, temp)
		}
	}
	ps.getRequestParameters(op, route)

	ps.getResponseParameters(op, route)
}

func (ps *Parser) getRequestParameters(op *openapi.Operation, route *spec.Route) {
	//process request body or form
	request := route.RequestType
	if request != nil {
		defineStruct, ok := request.(spec.DefineStruct)
		if ok {
			isBody := false
			members := defineStruct.Members
			for _, m := range members {
				isRequired := true
				tags := m.Tags()
				for _, tag := range tags {
					isRequired = !stringx.Contains(tag.Options, "optional")
					tempSchema := openapi.Schema{}
					ps.getParamType(&tempSchema, m)
					paramType := tempSchema.Type[0]
					//process header
					if tag.Key == "header" {
						op.Parameters = append(op.Parameters, openapi.Parameter{
							SimpleSchema: openapi.SimpleSchema{Type: paramType},
							ParamProps: openapi.ParamProps{
								Required:    isRequired,
								In:          InTypeMap["header"],
								Name:        tag.Name,
								Description: strings.ReplaceAll(m.Comment, "//", ""),
							},
						})
					}

					//process path
					if tag.Key == "path" {
						op.Parameters = append(op.Parameters, openapi.Parameter{
							SimpleSchema: openapi.SimpleSchema{Type: paramType},
							ParamProps: openapi.ParamProps{
								Required:    true,
								In:          InTypeMap["path"],
								Name:        tag.Name,
								Description: strings.ReplaceAll(m.Comment, "//", ""),
							},
						})
					}

					//process form
					if tag.Key == "form" {
						//if request is form and method is get, then field is query
						if strings.ToUpper(route.Method) == http.MethodGet {
							op.Parameters = append(op.Parameters, openapi.Parameter{
								SimpleSchema: openapi.SimpleSchema{Type: paramType},
								ParamProps: openapi.ParamProps{
									Required:    isRequired,
									In:          InTypeMap["getform"],
									Name:        tag.Name,
									Description: strings.ReplaceAll(m.Comment, "//", ""),
								},
							})
							continue
						}

						//post file
						if stringx.Contains(tag.Options, "file") {
							op.Consumes = openapi.StringOrArray{MineTypeMap["mpfd"]}
							op.Parameters = append(op.Parameters, openapi.Parameter{
								SimpleSchema: openapi.SimpleSchema{Type: "file"},
								ParamProps: openapi.ParamProps{
									Required:    isRequired,
									In:          InTypeMap["mpfd"],
									Name:        tag.Name,
									Description: strings.ReplaceAll(m.Comment, "//", ""),
								},
							})
							continue
						}

						//normal form request
						op.Consumes = openapi.StringOrArray{MineTypeMap["form"]}
						op.Parameters = append(op.Parameters, openapi.Parameter{
							SimpleSchema: openapi.SimpleSchema{Type: paramType},
							ParamProps: openapi.ParamProps{
								Required:    isRequired,
								In:          InTypeMap["postform"],
								Name:        tag.Name,
								Description: strings.ReplaceAll(m.Comment, "//", ""),
							},
						})
					}
				}

				if m.IsBodyMember() {
					isBody = true
					break
				}
			}
			//if request is body
			if isBody {
				op.Consumes = openapi.StringOrArray{MineTypeMap["json"]}
				description := strings.Join(request.Comments(), " ")
				if description == "" {
					description = "  "
				}
				ps.Swagger.Definitions[defineStruct.RawName] = ps.getDefinition(defineStruct)
				op.Parameters = append(op.Parameters, openapi.Parameter{
					ParamProps: openapi.ParamProps{
						Required:    isBody,
						In:          "body",
						Name:        request.Name(),
						Description: description,
						Schema: &openapi.Schema{
							SchemaProps: openapi.SchemaProps{
								Ref: openapi.MustCreateRef("#/definitions/" + defineStruct.RawName),
							},
						},
					},
				})
			}
		}
	}
}

func (ps *Parser) getResponseParameters(op *openapi.Operation, route *spec.Route) {
	response := route.ResponseType
	if response != nil {
		defineStruct, ok := response.(spec.DefineStruct)
		if ok {
			op.Responses = &openapi.Responses{
				ResponsesProps: openapi.ResponsesProps{
					StatusCodeResponses: make(map[int]openapi.Response),
				},
			}
			op.Responses.StatusCodeResponses[200] = openapi.Response{
				ResponseProps: openapi.ResponseProps{
					Description: "ok",
					Schema: &openapi.Schema{
						SchemaProps: openapi.SchemaProps{
							Ref: openapi.MustCreateRef("#/definitions/" + defineStruct.RawName),
						},
					},
				},
			}
		}
	}

	if len(op.Produces) < 1 {
		for _, p := range ps.Swagger.Produces {
			op.Produces = append(op.Produces, p)
		}
	}
}

func (ps *Parser) getParamType(schema *openapi.Schema, member spec.Member) {
	switch t := member.Type.(type) {
	case spec.PrimitiveType:
		schema.Type = openapi.StringOrArray{t.RawName}
	case spec.ArrayType:
		child := openapi.Schema{}
		ps.getParamType(&child, spec.Member{Type: t.Value})

		schema.Type = openapi.StringOrArray{"array"}
		schema.Items = &openapi.SchemaOrArray{}
		schema.Items.Schema = &child
	case spec.PointerType:
		ps.getParamType(schema, spec.Member{Type: t.Type})
	case spec.DefineStruct, spec.MapType, spec.InterfaceType:
		schema.Type = openapi.StringOrArray{"object"}
		schema.Ref = openapi.MustCreateRef("#/definitions/" + t.Name())
	default:
		panic("unknown type")
	}
}
