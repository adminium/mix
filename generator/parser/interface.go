package parser

import (
	"fmt"
	"go/ast"
)

type Interface struct {
	Name          string
	Methods       []*Method
	Defs          []*Def
	defs          map[string]*Def
	Imports       []*Import
	imports       map[string]bool
	Includes      []*Interface
	interfaceType *ast.InterfaceType
	file          *File
	loaded        bool
}

func (i *Interface) Load() (err error) {
	if i.loaded {
		return
	}
	for _, m := range i.interfaceType.Methods.List {
		switch mt := m.Type.(type) {
		case *ast.FuncType:
			cmt := i.file.comments.Filter(m).Comments()
			i.Methods = append(i.Methods, i.parseMethod(m.Names[0].Name, mt, cmt))
		case *ast.SelectorExpr:
			pkgName := mt.X.(*ast.Ident).String()
			typeName := mt.Sel.Name
			imt := i.file.getImport(pkgName)
			if imt == nil {
				panic(i.file.Errorf(m.Pos(), "can't found import: %s", pkgName))
			}
			include := imt.Package.GetInterface(typeName)
			if include == nil {
				panic(i.file.Errorf(m.Pos(), "can't found interface: %s.%s", pkgName, typeName))
			}
			err = include.Load()
			if err != nil {
				panic(err)
			}
			i.Includes = append(i.Includes, include)
			i.Methods = append(i.Methods, include.Methods...)
		case *ast.Ident:
			include := i.file.pkg.GetInterface(mt.Name)
			if include == nil {
				panic(i.file.Errorf(m.Pos(), "can't found interface: %s", mt.Name))
			}
			err = include.Load()
			if err != nil {
				panic(err)
			}
			i.Includes = append(i.Includes, include)
			i.Methods = append(i.Methods, include.Methods...)
		default:
			panic(i.file.Errorf(m.Pos(), "unsupported interface embed type"))
		}
	}
	i.loaded = true
	return
}

func (i *Interface) parseMethod(name string, t *ast.FuncType, cmt []*ast.CommentGroup) (r *Method) {
	
	r = &Method{Name: name}
	
	if t.Params != nil {
		for index, f := range t.Params.List {
			names := i.parseNames(f.Names)
			if len(names) == 0 {
				names = append(names, fmt.Sprintf("p%d", index))
			}
			r.Params = append(r.Params, i.parseParam(names, f.Type))
		}
	}
	
	if t.Results != nil {
		for index, f := range t.Results.List {
			names := i.parseNames(f.Names)
			if len(names) == 0 {
				names = append(names, fmt.Sprintf("r%d", index))
			}
			r.Results = append(r.Results, i.parseParam(names, f.Type))
		}
	}
	
	return
}

func (i *Interface) addDef(t *Def) {
	if i.defs == nil {
		i.defs = map[string]*Def{}
	}
	if _, ok := i.defs[t.Name]; ok {
		return
	}
	i.defs[t.Name] = t
	i.Defs = append(i.Defs, t)
}

func (i *Interface) getDef(name string) *Def {
	if i.defs == nil {
		return nil
	}
	return i.defs[name]
}

func (i *Interface) parseParam(names []string, t ast.Expr) (r *Param) {
	
	r = &Param{
		Names: names,
	}
	
	r.Type = parseType(i.file, i, "", t)
	
	return
}

func (i *Interface) addImport(imt *Import) {
	if i.imports == nil {
		i.imports = map[string]bool{}
	}
	if _, ok := i.imports[imt.Path]; ok {
		return
	}
	i.imports[imt.Path] = true
	i.Imports = append(i.Imports, imt)
}

func (i Interface) parseNames(idents []*ast.Ident) []string {
	names := make([]string, 0)
	for _, i := range idents {
		names = append(names, i.Name)
	}
	return names
}
