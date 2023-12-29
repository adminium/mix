package openapi

import (
	"fmt"
	"os"
	"strings"

	"github.com/adminium/fs"
	"github.com/adminium/mix/generator/parser"
	"github.com/adminium/pointer"
	"github.com/gozelle/openapi/openapi3"
	"github.com/invopop/yaml"
)

type GenerateCmd struct {
	Path      string
	Interface string
}

func Parse(tplFile, sourcePath string, interfaceName string) (doc *DocumentV3, err error) {

	doc, err = Load(tplFile)
	if err != nil {
		err = fmt.Errorf("load template error: %s", err)
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	mod, err := parser.PrepareMod(pwd)
	if err != nil {
		err = fmt.Errorf("prepare mod error: %s", err)
		return
	}

	pkg, err := parser.Parse(mod, sourcePath)
	if err != nil {
		err = fmt.Errorf("parse package error: %s", err)
		return
	}

	i := pkg.GetInterface(interfaceName)
	if i == nil {
		err = fmt.Errorf("api interface: %s not found", interfaceName)
		return
	}

	err = i.Load()
	if err != nil {
		err = fmt.Errorf("load interface: %s error: %s", interfaceName, err)
		return
	}

	r := ConvertAPI(i)
	ConvertOpenapi(doc, r)

	return
}

func Load(file string) (doc *DocumentV3, err error) {
	doc = &DocumentV3{}
	if file != "" {
		var c []byte
		c, err = fs.Read(file)
		if err != nil {
			err = fmt.Errorf("read openapi file error: %s", err)
			return
		}
		if strings.HasSuffix(file, ".json") {
			err = doc.UnmarshalJSON(c)
			if err != nil {
				err = fmt.Errorf("unmarshal openapi file error: %s", err)
				return
			}
		} else if strings.HasSuffix(file, ".yaml") {
			var j []byte
			j, err = yaml.YAMLToJSON(c)
			if err != nil {
				err = fmt.Errorf("convert yaml to json error: %s", err)
				return
			}

			err = doc.UnmarshalJSON(j)
			if err != nil {
				err = fmt.Errorf("unmarshal openapi file error: %s", err)
				return
			}
		} else {
			err = fmt.Errorf("unsupport openapi file: %s suffix, accept json or yaml", file)
			return
		}
	}

	return
}

func ConvertOpenapi(doc *DocumentV3, r *API) {

	doc.OpenAPI = "3.0.3"

	if doc.Info == nil {
		doc.Info = &openapi3.Info{}
	}

	for _, v := range r.Methods {
		convertOpenapiMethods(doc, v)
	}

	if doc.Components == nil {
		doc.Components = &openapi3.Components{}
	}
	if doc.Components.Schemas == nil {
		doc.Components.Schemas = map[string]*openapi3.SchemaRef{}
	}

	for _, v := range r.Defs {
		if ref := makeOpenapiSchemaRef(doc, v); ref != nil {
			doc.Components.Schemas[v.Name] = ref
		}
	}

	return
}

func convertOpenapiMethods(d *DocumentV3, m *Method) {
	if d.Paths == nil {
		d.Paths = map[string]*openapi3.PathItem{}
	}
	item := &openapi3.PathItem{
		Summary:     "",
		Description: "",
		Post: &openapi3.Operation{
			Tags:        nil,
			Summary:     "",
			Description: "",
			OperationID: m.Name,
			Parameters:  nil,
			RequestBody: nil,
			Responses:   nil,
		},
	}
	if m.Request != nil {
		if m.Request.Type == ArrayParams {
			for _, v := range m.Request.ArrayFields {
				ref := makeOpenapiSchemaRef(d, v)
				if ref == nil {
					panic(fmt.Errorf("parse param: %s get nil ref", v.Field))
				}
				param := &openapi3.ParameterRef{Value: &openapi3.Parameter{}}
				param.Value.Name = v.Field
				param.Value.Schema = ref
				item.Post.Parameters = append(item.Post.Parameters, param)
			}
		} else {
			if ref := makeOpenapiMethodParameterRef(d, m.Request); ref != "" {
				item.Post.RequestBody = &openapi3.RequestBodyRef{
					Ref: ref,
				}
			}
		}
	}
	if item.Post.Responses == nil {
		item.Post.Responses = map[string]*openapi3.ResponseRef{}
	}
	// 过滤响应 io.Reader 和 chan
	if m.Replay != nil {
		if m.Replay.Use != nil {
			if m.Replay.Use.Type == parser.TStruct {
				item.Post.Responses["200"] = &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: pointer.ToString("success"),
						Content: openapi3.Content{
							ApplicationJson: &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Ref: fmt.Sprintf("#/components/schemas/%s", m.Replay.Use.Name),
								},
							},
						},
					},
				}
			} else {
				item.Post.Responses["200"] = &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: pointer.ToString("success"),
						Content:     makeOpenapiContent(d, m.Replay.Use),
					},
				}
			}
		} else {
			item.Post.Responses["200"] = &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: pointer.ToString("success"),
					Content:     makeOpenapiContent(d, m.Replay),
				},
			}
		}
	}

	if item.Post.Responses["200"] == nil {
		item.Post.Responses["200"] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: pointer.ToString("success"),
			},
		}
	}

	d.Paths[fmt.Sprintf("/%s", m.Name)] = item
}

func makeOpenapiMethodParameterRef(d *DocumentV3, def *Def) string {
	if d.Components == nil {
		d.Components = &openapi3.Components{}
	}
	if d.Components.RequestBodies == nil {
		d.Components.RequestBodies = map[string]*openapi3.RequestBodyRef{}
	}
	ref := fmt.Sprintf("#/components/requestBodies/%s", def.Field)
	if def.Use == nil {
		d.Components.RequestBodies[def.Field] = &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: false,
				Content:  makeOpenapiContent(d, def),
			},
		}
		return ref
	} else if def.Use.Type == parser.TStruct {
		d.Components.RequestBodies[def.Field] = &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: false,
				Content: openapi3.Content{
					ApplicationJson: &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: fmt.Sprintf("#/components/schemas/%s", def.Use.Name),
						},
					},
				},
			},
		}
		return ref
	}

	return ""
}

func makeOpenapiContent(d *DocumentV3, def *Def) openapi3.Content {
	var c openapi3.Content = map[string]*openapi3.MediaType{}
	if ref := makeOpenapiSchemaRef(d, def); ref != nil {
		c[ApplicationJson] = &openapi3.MediaType{
			Schema: ref,
		}
	}
	return c
}

func makeOpenapiSchemaRef(d *DocumentV3, def *Def) (s *openapi3.SchemaRef) {

	s = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: convertOpenapiType(def),
		},
	}

	if def.Use != nil && def.Use.Type == parser.TStruct {
		s.Ref = fmt.Sprintf("#/components/schemas/%s", def.Use.Name)
	} else if def.Type == parser.TStruct {
		s.Value.Properties = map[string]*openapi3.SchemaRef{}
		for _, v := range def.StructFields {
			name := v.Field
			if v.Json != "" {
				name = v.Json
			}
			if ref := makeOpenapiSchemaRef(d, v); ref != nil {
				s.Value.Properties[name] = ref
			}
		}
		if len(s.Value.Properties) == 0 {
			return nil
		}
	} else if def.Type == parser.TSlice {
		if ref := makeOpenapiSchemaRef(d, def.Elem); ref != nil {
			s.Value.Items = ref
		}
	}
	return
}

func convertOpenapiType(def *Def) string {

	typ := def.Type
	if def.Use != nil {
		typ = def.Use.Type
	}

	switch typ {
	case parser.TSlice, parser.TArray:
		return Array
	case parser.TString:
		return String
	case parser.TInt, parser.TInt8, parser.TInt16, parser.TInt32, parser.TInt64,
		parser.TUint, parser.TUint8, parser.TUint16, parser.TUint32, parser.TUint64:
		return Integer
	case parser.TFloat32, parser.TFloat64:
		return Number
	case parser.TBool:
		return Boolean
	}

	return Object
}
