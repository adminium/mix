package tests_basic

import (
	"context"
)

type BasicAPI interface {
	Test1(ctx context.Context, in Basic) (out Basic, err error)
	Test2(context.Context, *Basic) (*Basic, error)
}

type Basic struct {
	Int              int
	Uint             *uint
	Bool             bool
	Float            *float32
	String           string
	Map              map[string]interface{}
	IntArray         []int      `json:"int_array"`
	UintArray        []*uint    `json:"uint_array"`
	BoolArray        []bool     `json:"bool_array"`
	FloatArray       []*float32 `json:"float_array"`
	StringArray      []string   `json:"string_array"`
	StringArrayArray [][]string `json:"string_array_array"`
}
