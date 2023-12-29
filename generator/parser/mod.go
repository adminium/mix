package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/adminium/fs"
	"github.com/gozelle/color"
	"golang.org/x/mod/modfile"
)

// PrepareMod prepare Mod object in the specified directory
func PrepareMod(dir string) (mod *Mod, err error) {

	defer func() {
		if err == nil {
			if mod.file.Module == nil {
				err = fmt.Errorf("no module is defined in the mod file")
				return
			}
		}
	}()

	for {
		var f *os.File
		if f, err = os.Open(fs.Join(dir, "go.mod")); err == nil {
			var d []byte
			d, err = io.ReadAll(f)
			if err != nil {
				return
			}
			_ = f.Close()
			if err == nil || err == io.EOF {
				// go.mod exists and is readable (is a file, not a directory).
				var mf *modfile.File
				mf, err = modfile.Parse("go.mod", d, nil)
				if err != nil {
					return
				}
				mod = &Mod{file: mf, root: dir}
				break
			}
		}
		d := filepath.Dir(dir)
		if len(d) >= len(dir) {
			return
		}
		dir = fs.Join(dir, "../")
	}
	return
}

type Mod struct {
	root     string
	file     *modfile.File
	packages map[string]*Package
	loaded   map[string]bool
}

func (m Mod) ModuleName() string {
	if m.file == nil {
		return ""
	}
	return m.file.Module.Mod.Path
}

func (m Mod) Gopath() string {
	return os.Getenv("GOPATH")
}

func (m *Mod) cachePackage(realPath string, pkg *Package) {
	if m.packages == nil {
		m.packages = map[string]*Package{}
	}

	m.packages[realPath] = pkg
}

// GetPackageLocation 获取包的真实路径
// 1. 首先判断是否被本地替换
// 2. 然后判断是否是直接依赖的包
// 3. 最后判断是否为系统包
func (m Mod) GetPackageLocation(pkg string) string {
	if m.file == nil {
		return ""
	}
	if m.ModuleName() != "std" {
		for _, v := range m.file.Require {
			if v.Mod.Path == pkg {
				return fs.Join(m.Gopath(), "pkg/mod", fmt.Sprintf("%s@%s", pkg, v.Mod.Version))
			} else {
				c := pkg
				valid := false
				for c != "." && c != "/" {
					if v.Mod.Path == c {
						valid = true
						break
					}
					c = filepath.Join(c, "../")
				}
				if valid {
					// 如果有根项目定义包版本升级
					for _, pv := range m.file.Require {
						if pv.Mod.Path == c {
							return fs.Join(m.Gopath(), "pkg/mod", fmt.Sprintf("%s@%s%s", c, pv.Mod.Version, strings.TrimPrefix(pkg, c)))
						}
					}
					return fs.Join(m.Gopath(), "pkg/mod", fmt.Sprintf("%s@%s%s", c, v.Mod.Version, strings.TrimPrefix(pkg, c)))
				}
			}
		}
		for _, v := range m.file.Replace {
			if v.New.Path == pkg {
				return v.Old.Path
			} else {
				c := pkg
				valid := false
				for c != "." && c != "/" {
					if v.New.Path == c {
						valid = true
						break
					}
					c = filepath.Join(c, "../")
				}
				if valid {
					return fmt.Sprintf("%s%s", v.Old.Path, strings.TrimPrefix(pkg, c))
				}
			}
		}
	}

	if !fs.Exist(fs.Join(m.Gopath(), "src")) {
		fmt.Println(color.YellowString("no src directory found at GOPATH: %s\n", m.Gopath()))
	}

	path := fs.Join(m.Gopath(), "src", pkg)
	if fs.Exist(path) {
		return path
	}

	path = fs.Join(m.Gopath(), "src/vendor", pkg)
	if fs.Exist(path) {
		return path
	}

	path = fs.Join(m.root, strings.TrimPrefix(strings.TrimPrefix(pkg, m.file.Module.Mod.Path), "/"))
	if fs.Exist(path) {
		return path
	}

	return ""
}

func (m Mod) GetPackageRealName(path string) (name string, err error) {

	files, err := m.GetPackageFiles(path)
	if err != nil {
		return
	}
	set := token.NewFileSet()

	var f *ast.File
	for _, v := range files {
		if strings.HasSuffix(v, "_test.go") {
			continue
		}

		f, err = parser.ParseFile(set, v, nil, parser.AllErrors|parser.ParseComments)
		if err != nil {
			return
		}
		if f.Name == nil {
			err = fmt.Errorf("package name is nil")
			return
		}
		name = f.Name.String()

		// ignore main namespace
		if name != "main" {
			return
		}
	}

	return
}

func (m Mod) GetPackageFiles(path string) (files []string, err error) {
	if !fs.Exist(path) {
		return
	}
	files, err = fs.Files(false, path, ".go")
	if err != nil {
		return
	}
	return
}
