package parser

import (
	"encoding/json"
	"go/ast"
)

type Def struct {
	Name     string   `json:"Name"`
	Expr     ast.Expr `json:"-"`
	File     *File    `json:"-"`
	Type     *Type    `json:"Type,omitempty"`
	ToString bool     `json:"ToString,omitempty"`
	IsStrut  bool     `json:"IsStrut,omitempty"`
	parsed   bool
}

func (d Def) ShallowFork() *Def {
	n := &Def{
		Name:     d.Name,
		Expr:     d.Expr,
		File:     d.File,
		Type:     d.Type,
		ToString: d.ToString,
		IsStrut:  d.IsStrut,
		parsed:   d.parsed,
	}
	return n
}

func (d Def) String() string {
	a, _ := json.Marshal(d)
	return string(a)
}
