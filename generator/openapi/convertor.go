package openapi

import (
	"fmt"

	"github.com/adminium/logger"
	"github.com/adminium/mix/generator/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var log = logger.NewLogger("convertor")

func ConvertAPI(i *parser.Interface) *API {

	r := &API{
		Name: i.Name,
	}

	//for _, v := range i.Imports {
	//	r.Imports = append(r.Imports, &Import{
	//		Alias: v.Alias,
	//		Path:  v.Path,
	//	})
	//}

	convertDefs(r, i.Defs)
	for _, v := range i.Inlines {
		convertDefs(r, v.Defs)
	}

	for _, v := range i.Methods {
		r.Methods = append(r.Methods, convertRenderMethod(v))
	}

	return r
}

func convertDefs(r *API, defs []*parser.Def) {
	for _, v := range defs {
		if v.Type.Type != parser.TStruct {
			continue
		}
		d := convertRenderDef(v)
		r.Defs = append(r.Defs, d)
	}
}

func convertRenderMethod(m *parser.Method) *Method {
	r := &Method{
		Name:    m.Name,
		Request: convertRenderMethodRequest(m),
		Replay:  convertRenderMethodReply(m),
	}
	return r
}

func convertRenderMethodRequest(m *parser.Method) *Def {
	request := &Def{
		Type: ArrayParams,
	}
	params := m.ExportParams()
	if len(params) == 1 && params[0].Type.NoPointer().Def != nil && params[0].Type.NoPointer().Def.Type.RealType().IsStruct() {
		request.ArrayFields = append(request.ArrayFields, convertRenderMethodParam(params[0])[0])
	} else if len(params) > 0 {
		for _, v := range params {
			request.ArrayFields = append(request.ArrayFields, convertRenderMethodParam(v)...)
		}
	} else {
		request = nil
	}

	return request
}

func convertRenderMethodReply(m *parser.Method) *Def {
	var replay *Def

	results := m.ExportResults()

	if len(results) == 0 {
		return nil
	}

	if results[0].Type.NoPointer().Def != nil && results[0].Type.NoPointer().Def.Type.RealType().IsStruct() {
		replay = &Def{Use: convertRenderDef(results[0].Type.NoPointer().Def)}
	} else {
		replay = convertRenderType(results[0].Type)
	}
	replay.Name = fmt.Sprintf("%sReply", m.Name)
	return replay
}

func convertRenderMethodParam(p *parser.Param) []*Def {
	r := make([]*Def, 0)
	for _, v := range p.Names {
		d := convertRenderType(p.Type)
		d.Field = v
		d.Json = v
		r = append(r, d)
	}
	return r
}

func Title(v string) string {
	return cases.Title(language.English).String(v)
}

func convertRenderDef(d *parser.Def) *Def {
	n := convertRenderType(d.Type)
	n.Name = d.Name
	return n
}

func convertRenderType(t *parser.Type) *Def {

	pointer := t.Pointer
	rt := t.NoPointer().RealType()

	n := &Def{
		Field:        t.Field,
		Json:         rt.Json(),
		Type:         rt.Type,
		Pointer:      pointer,
		StructFields: nil,
		Elem:         nil,
		Use:          nil,
		Tags:         string(rt.Tags),
	}

	switch rt.Type {
	case parser.TStruct:
		for _, v := range rt.StructFields {
			n.StructFields = append(n.StructFields, convertRenderType(v))
		}
	case parser.TSlice, parser.TArray:
		n.Type = parser.TSlice
		n.Elem = convertRenderType(rt.Elem)
	}

	// 处理引用关系
	if rt.Def != nil {
		n.Use = &Def{Name: rt.Def.Name}
		if rt.Def.IsStrut {
			n.Use.Type = parser.TStruct
		} else if rt.Def.Type != nil {
			n.Use.Type = rt.Def.Type.Type
		}
	}

	return n
}
