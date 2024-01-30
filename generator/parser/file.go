package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

type File struct {
	path     string
	mod      *Mod
	pkg      *Package
	imports  map[string]*Import
	comments ast.CommentMap
	ast      *ast.File
	set      *token.FileSet
}

func (f File) Errorf(pos token.Pos, format string, a ...any) error {
	p := f.set.Position(pos)
	return fmt.Errorf("parse file: %s error: %s", p, fmt.Sprintf(format, a...))
}

func (f *File) getImport(name string) *Import {
	if f.imports == nil {
		return nil
	}
	i := f.imports[name]
	if i != nil {
		f.loadImport(i)
	}
	return i
}

func (f *File) addImport(item *Import) {
	if f.imports == nil {
		f.imports = map[string]*Import{}
	}

	f.imports[item.Name] = item
}

func (f *File) markToStringDef(node ast.Node) {
	d, ok := node.(*ast.FuncDecl)
	if !ok {
		return
	}
	if d.Recv == nil {
		return
	}
	if d.Recv.NumFields() == 0 {
		return
	}

	if d.Type.Results == nil {
		return
	}
	if d.Type.Results.NumFields() != 1 {
		return
	}
	s, ok := d.Type.Results.List[0].Type.(*ast.Ident)
	if !ok {
		return
	}
	if s.Name != "string" {
		return
	}

	i, ok := d.Recv.List[0].Type.(*ast.Ident)
	if !ok {
		return
	}

	f.pkg.markStringer(i.String())
}

func (f *File) Visit(node ast.Node) ast.Visitor {

	f.markToStringDef(node)

	s, ok := node.(*ast.TypeSpec)
	if !ok {
		return f
	}
	switch t := s.Type.(type) {
	case *ast.InterfaceType:
		f.pkg.addInterface(s.Name.String(), &Interface{Name: s.Name.String(), astInterface: t, file: f})
	}

	switch t := s.Type.(type) {

	case *ast.InterfaceType,
		*ast.FuncType,
		*ast.StructType,
		*ast.Ident,
		*ast.MapType,
		*ast.SliceExpr,
		*ast.ArrayType,
		*ast.StarExpr,
		*ast.SelectorExpr,
		*ast.ChanType,
		*ast.IndexExpr,
		*ast.IndexListExpr:
		if !s.Name.IsExported() {
			return f
		}
		_, isStruct := t.(*ast.StructType)
		f.pkg.addDef(s.Name.String(), &Def{Name: s.Name.String(), File: f, Expr: s.Type, IsStrut: isStruct})

	default:
		panic(fmt.Errorf("unsupport parse type: %s at: %s", reflect.TypeOf(t), f.set.Position(s.Pos())))
	}

	return f
}

func (f *File) parse(file string) (err error) {
	set := token.NewFileSet()
	af, err := parser.ParseFile(set, file, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return
	}
	f.ast = af
	f.set = set
	f.comments = ast.NewCommentMap(set, af, af.Comments)

	// cache the import declarations for the current file,
	// with the identifier taking an importing alias or package path
	for _, i := range af.Imports {
		v := f.prepareImport(i)
		f.addImport(v)
	}

	// parse file use Visit
	ast.Walk(f, af)

	return
}

func (f *File) parseAlias(name string) string {
	if name == "<nil>" {
		return ""
	}
	return name
}

func (f File) prepareImport(i *ast.ImportSpec) *Import {
	r := &Import{
		Name: f.parseAlias(i.Name.String()),
		Package: &Package{
			Path: strings.Trim(i.Path.Value, "\""),
		},
	}
	// 忽略 go.mod 不依赖的包
	location := f.mod.GetPackageLocation(r.Path)
	if location == "" {
		panic(fmt.Errorf("can't resolve package: %s location", r.Path))
	}
	r.location = location

	var err error
	if r.Name == "" {
		r.Name, err = f.mod.GetPackageRealName(location)
		if err != nil {
			panic(fmt.Errorf("%s: get %s package name error: %s", f.path, r.Path, err))
		}
	}
	return r
}

func (f *File) loadImport(r *Import) {

	if f.mod.packages != nil {
		if v, ok := f.mod.packages[r.Path]; ok {
			r.Package = v
			return
		}
	}

	//if strings.Contains(r.Path, "/internal") {
	//	return nil
	//}

	if r.Path == "C" {
		return
	}

	var err error

	defer func() {
		f.mod.cachePackage(r.Path, r.Package)
	}()

	err = r.Package.Parse(f.mod, r.location)
	if err != nil {
		panic(fmt.Errorf("load package: %s error: %s", r.Path, err))
	}
	return
}
