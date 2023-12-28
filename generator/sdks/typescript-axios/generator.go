package typescript_axios

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/adminium/mix/generator/openapi"
	"github.com/adminium/mix/generator/writter"
	"github.com/fatih/structs"
	"github.com/flosch/pongo2/v6"
	"github.com/gozelle/openapi/openapi3"
)

type API struct {
	Name    string    `json:"Name,omitempty"`
	Methods []*Method `json:"Methods,omitempty"`
	Types   []*Type   `json:"Types,omitempty"`
}

type Method struct {
	Method   string `json:"Method,omitempty"`
	Path     string `json:"Path,omitempty"`
	Name     string `json:"Name,omitempty"`
	Request  string `json:"Request,omitempty"`
	Response string `json:"Response,omitempty"`
	Data     string `json:"Data,omitempty"`
	Comment  string `json:"Comment,omitempty"`
}

type Type struct {
	Name   string  `json:"Name,omitempty"`
	Type   string  `json:"Type,omitempty"`
	Fields []*Type `json:"Fields,omitempty"`
	Elem   *Type   `json:"Elem,omitempty"`
	Ref    string  `json:"Ref,omitempty"`
}

type Option struct {
	Name string `json:"name"`
}

func Generate(file string, option string) (files []*writter.File, err error) {

	opt := Option{}
	if option != "" {
		err = json.Unmarshal([]byte(option), &opt)
		if err != nil {
			err = fmt.Errorf("parse option: '%s' error: %s", option, err)
			return
		}
	}
	if opt.Name == "" {
		opt.Name = "SDK"
	}

	doc, err := openapi.Load(file)
	if err != nil {
		err = fmt.Errorf("load openapi file error: %s", err)
		return
	}

	api, err := Convert(doc)
	if err != nil {
		err = fmt.Errorf("convert ts type error: %s", err)
		return
	}

	//fmt.Println("print API:")
	//spew.Json(api)

	tpl, err := pongo2.FromString(apiTpl)
	if err != nil {
		err = fmt.Errorf("preprea tpl error: %s", err)
		return
	}

	params := structs.Map(api)
	params["opt"] = opt

	c, err := tpl.Execute(params)
	files = append(files, []*writter.File{
		{Name: fmt.Sprintf("%s.ts", opt.Name), Content: c},
		{Name: "base.ts", Content: baseTpl},
	}...)

	return
}

func Convert(doc *openapi.DocumentV3) (api *API, err error) {

	api = &API{
		Name: "SDK",
	}

	for path, item := range doc.Paths {
		var m *Method
		m, err = convertMethod(doc, path, item)
		if err != nil {
			return
		}
		api.Methods = append(api.Methods, m)
	}

	if doc.Components != nil {
		if doc.Components.Schemas != nil {
			api.Types = append(api.Types, convertSchemas(doc, doc.Components.Schemas)...)
		}
		if doc.Components.RequestBodies != nil {
			api.Types = append(api.Types, convertRequestBodies(doc, doc.Components.RequestBodies)...)
		}
		if doc.Components.Responses != nil {
			api.Types = append(api.Types, convertResponses(doc, doc.Components.Responses)...)
		}
	}

	return
}

func convertMethod(doc *openapi.DocumentV3, path string, item *openapi3.PathItem) (m *Method, err error) {
	var op *openapi3.Operation
	var method string
	if item.Post != nil {
		op = item.Post
		method = "post"
	} else if item.Get != nil {
		op = item.Get
		method = "get"
	} else if item.Put != nil {
		op = item.Put
		method = "put"
	} else if item.Delete != nil {
		op = item.Delete
		method = "delete"
	} else {
		err = fmt.Errorf("handle path: %s error: only accept POST、GET、PUT、DELETE method", path)
		return
	}

	m = &Method{
		Method:   method,
		Path:     path,
		Name:     op.OperationID,
		Comment:  op.Description,
		Request:  "",
		Response: "",
	}

	m.Request, m.Data = convertMethodRequestBody(doc, op)

	if op.Responses != nil {
		m.Response = convertMethodResponses(doc, op.Responses)
	}

	return
}

func convertMethodRequestBody(doc *openapi.DocumentV3, op *openapi3.Operation) (t, d string) {

	if len(op.Parameters) > 0 {
		d = "["
		for _, v := range op.Parameters {
			if v.Value.Schema.Ref != "" {
				t += fmt.Sprintf("%s: %s,", v.Value.Name, filepath.Base(v.Value.Schema.Ref))

			} else {
				t += fmt.Sprintf("%s: %s,", v.Value.Name, convertSchema(doc, "", v.Value.Schema.Value).Type)
			}
			d += fmt.Sprintf("%s,", v.Value.Name)
		}
		t = strings.TrimSuffix(t, ",")
		d = strings.TrimSuffix(d, ",")
		d += "]"
		return
	}

	return
}

func convertMethodResponses(doc *openapi.DocumentV3, resp openapi3.Responses) (t string) {
	if resp["200"] != nil &&
		resp["200"].Value != nil &&
		resp["200"].Value.Content != nil &&
		resp["200"].Value.Content[openapi.ApplicationJson] != nil {
		if resp["200"].Value.Content[openapi.ApplicationJson].Schema != nil &&
			resp["200"].Value.Content[openapi.ApplicationJson].Schema.Ref != "" {
			t = filepath.Base(resp["200"].Value.Content[openapi.ApplicationJson].Schema.Ref)
		} else if resp["200"].Value.Content[openapi.ApplicationJson].Schema.Value != nil {
			t = convertSchema(doc, "", resp["200"].Value.Content[openapi.ApplicationJson].Schema.Value).Type
		}
	}
	return
}

func convertSchemas(doc *openapi.DocumentV3, schemas openapi3.Schemas) (types []*Type) {

	for name, v := range schemas {
		if vv := convertSchemaRef(doc, name, v); vv != nil {
			types = append(types, vv)
		}
	}

	return
}

func convertRequestBodies(doc *openapi.DocumentV3, reqs openapi3.RequestBodies) (types []*Type) {

	for name, v := range reqs {
		if v.Value != nil &&
			v.Value.Content != nil &&
			v.Value.Content[openapi.ApplicationJson] != nil &&
			v.Value.Content[openapi.ApplicationJson].Schema != nil {
			types = append(types, convertSchemaRef(doc, name, v.Value.Content[openapi.ApplicationJson].Schema))
		}
	}

	return
}

func convertResponses(doc *openapi.DocumentV3, reps openapi3.Responses) (types []*Type) {
	for name, v := range reps {
		if v.Value != nil &&
			v.Value.Content != nil &&
			v.Value.Content[openapi.ApplicationJson] != nil &&
			v.Value.Content[openapi.ApplicationJson].Schema != nil {
			types = append(types, convertSchemaRef(doc, name, v.Value.Content[openapi.ApplicationJson].Schema))
		}
	}
	return
}

func convertSchemaRef(doc *openapi.DocumentV3, name string, schema *openapi3.SchemaRef) (t *Type) {
	if schema.Ref != "" {
		if strings.HasPrefix(schema.Ref, "#/components/schemas") {
			if v := getSchemaRef(doc, filepath.Base(schema.Ref)); v != nil {
				schema.Value = v.Value
			}

		}
	}

	if schema.Value != nil {
		t = convertSchema(doc, name, schema.Value)
	}

	return
}

func convertSchema(doc *openapi.DocumentV3, name string, value *openapi3.Schema) (field *Type) {
	field = &Type{
		Name: name,
	}
	switch value.Type {
	case openapi.Object:
		if value.Properties != nil {
			field.Type = openapi.Object
			for k, v := range value.Properties {
				if v.Ref != "" {
					field.Fields = append(field.Fields, &Type{Name: k, Type: filepath.Base(v.Ref)})
				} else if v.Value != nil {
					field.Fields = append(field.Fields, convertSchema(doc, k, v.Value))
				}
			}
			o, err := pongo2.FromString(objectTpl)
			if err != nil {
				panic(fmt.Errorf("prepare object tpl error: %s", err))
			}
			t, err := o.Execute(structs.Map(field))
			if err != nil {
				panic(fmt.Errorf("render object tpl error: %s", err))
			}
			field.Type = t
		} else {
			field.Type = "any"
		}
	case openapi.Array:
		typ := "any"
		if value.Items != nil {
			v := value.Items.Value
			if value.Items.Ref != "" {
				typ = filepath.Base(value.Items.Ref)
			} else if v != nil {
				typ = convertSchema(doc, "", v).Type
			}
		}
		field.Type = fmt.Sprintf("%s[]", typ)
	case openapi.String:
		field.Type = openapi.String
	case openapi.Boolean:
		field.Type = openapi.Boolean
	case openapi.Integer, openapi.Number:
		field.Type = openapi.Number
	}

	return
}

func getSchemaRef(doc *openapi.DocumentV3, name string) *openapi3.SchemaRef {
	if doc.Components != nil &&
		doc.Components.Schemas != nil {
		ref := doc.Components.Schemas[name]
		if ref != nil {
			if ref.Ref != "" {
				return getSchemaRef(doc, filepath.Base(ref.Ref)) // TODO 递归获取引用，可以加上计数避免死循环
			}
			return ref
		}
	}
	return nil
}

const apiTpl = `
/* tslint:disable */
/* eslint-disable */
/**
 *
 * No description provided (generated by Mix Openapi Generator https://github.com/gozelle/mix)
 *
 * The version of the OpenAPI document: 3.0.3
 *
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://mix.gozelle.io).
 * Do not edit the class manually.
 */
import {BaseAPI} from "./base";
// @ts-ignore
import {AxiosInstance,AxiosRequestConfig} from "axios";


/**
Example:
  export const api = new {{ opt.Name }}(axios.create({
    baseURL: 'http://127.0.0.1:8080/api/v1/module',
  }));
  api.client.interceptors.request.use(...requestInterceptorExample);
  api.client.interceptors.response.use(...responseInterceptorExample);
*/
export class {{ opt.Name }} extends BaseAPI {

	constructor(instance: AxiosInstance) {
        super(instance);
    }
{%- for method in Methods %}

    public {{method.Name}}({% if method.Request %}{{method.Request}}, {% endif %}options?: AxiosRequestConfig): Promise<{% if method.Response %}{{method.Response}}{% else %}null{% endif %}> {
        return this.client.{{ method.Method }}('{{ method.Path }}', {% if method.Request %}{{method.Data}}{% else %}null{% endif %}, options)
    }

{%- endfor %}
}

{%- for type in Types %}

export interface {{ type.Name }} {
	{%- for field in type.Fields %}
    {{ field.Name }}?: {{ field.Type }};
	{%- endfor %}
}
{%- endfor %}
`

const objectTpl = `
{
	{%- for field in Fields %}
	    {{ field.Name }}?: {{ field.Type }};
	{%- endfor %}
}
`
