package parser

import (
	"fmt"
	"strings"

	"github.com/adminium/fs"
)

type Packages map[string]*Package

type Package struct {
	Path       string
	interfaces map[string]*Interface
	defs       map[string]*Def
	defsCount  map[string]int
	stringers  map[string]bool
}

func (p *Package) Defs() map[string]*Def {
	return p.defs
}

func (p *Package) markStringer(t string) {
	if p.stringers == nil {
		p.stringers = map[string]bool{}
	}
	p.stringers[t] = true
}

func (p *Package) isStringer(t string) bool {
	if p.stringers == nil {
		return false
	}
	_, ok := p.stringers[t]
	return ok
}

func (p *Package) GetDef(name string) *Def {
	if p.defs == nil {
		return nil
	}
	return p.defs[name]
}

func (p *Package) addDef(name string, item *Def) {
	if p.defs == nil {
		p.defs = map[string]*Def{}
	}
	p.defs[name] = item
}

func (p *Package) addInterface(name string, item *Interface) *Interface {
	if p.interfaces == nil {
		p.interfaces = map[string]*Interface{}
	}
	if v, ok := p.interfaces[name]; ok {
		return v
	}
	p.interfaces[name] = item
	return item
}

func (p *Package) GetInterface(name string) *Interface {
	if p.interfaces == nil {
		return nil
	}
	v, ok := p.interfaces[name]
	if ok {
		return v
	}
	return nil
}

func (p *Package) parseFiles(mod *Mod, files []string) (err error) {
	for _, v := range files {
		if !strings.HasSuffix(v, "_test.go") {
			f := &File{pkg: p, mod: mod, path: v}
			err = f.parse(v)
			if err != nil {
				return
			}
		}
	}
	return
}

func (p *Package) Parse(mod *Mod, dir string) error {

	if mod.loaded == nil {
		mod.loaded = map[string]bool{}
	}
	if _, ok := mod.loaded[dir]; ok {
		return nil
	}

	mod.loaded[dir] = true
	if !fs.IsDir(dir) {
		return fmt.Errorf("only accept dir, get: %s", dir)
	}
	files, err := fs.Files(false, dir, ".go")
	if err != nil {
		return err
	}

	err = p.parseFiles(mod, files)
	if err != nil {
		return err
	}

	for _, v := range p.defs {
		if p.isStringer(v.Name) {
			v.ToString = true
		}
	}

	return nil
}

func (p *Package) AddExternalNalDef(def *Def) {
	if def.File.pkg.Path == "context" {
		return
	}
	d := def.ShallowFork()
	t := p.GetDef(d.Name)
	if t != nil {
		if p.defsCount == nil {
			p.defsCount = map[string]int{}
		}
		if v, ok := p.defsCount[d.Name]; !ok {
			p.defsCount[d.Name] = 1
			d.Name = fmt.Sprintf("%s%d", d.Name, 2)
		} else {
			d.Name = fmt.Sprintf("%s%d", d.Name, v+1)
			p.defsCount[d.Name] = v + 1
		}
	}
	p.addDef(d.Name, d)
}
