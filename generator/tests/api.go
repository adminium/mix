package tests

import (
	"context"
	"github.com/shopspring/decimal"
	"time"
)

type RequestAPI interface {
	Ping(ctx context.Context) (err error)
	Int(ctx context.Context, p1 int, p2 int8, p3 int16, p4 int32, p5 int64) error
	IntPointer(context.Context, *int, *int8, *int16, *int32, *int64) error
	IntArray(ctx context.Context, p1 []int) error
	IntArrayPointer(context.Context, []*int) error
	IntStruct(ctx context.Context, param1 Int) error
	IntStructPointer(context.Context, *IntPointer) error
	IntStructArray(ctx context.Context, p1 []Int) error
	IntStructPointerArray(context.Context, []*IntPointer) error
	Uint(ctx context.Context, p1 uint, p2 uint8, p3 uint16, p4 uint32, p5 uint64) error
	UIntPointer(context.Context, *uint, *uint8, *uint16, *uint32, *int64) error
	UintArray(ctx context.Context, p1 []uint) error
	UintArrayPointer(context.Context, []*uint) error
	UintStruct(ctx context.Context, p1 Uint) error
	UintStructPointer(context.Context, *UintPointer) error
	String(ctx context.Context, p1 string) error
	StringPointer(context.Context, *string) error
	StringArray(ctx context.Context, p1 []string) error
	StringArrayPointer(context.Context, []*string) error
	StringStruct(ctx context.Context, p1 String) error
	StringStructPointer(context.Context, *StringPointer) error
	Float(ctx context.Context, p1 float32, p2 float64) error
	FloatPointer(context.Context, *float32, *float64) error
	FloatArray(ctx context.Context, p1 []float32, p2 []float64) error
	FloatArrayPointer(context.Context, []*float32, []*float64) error
	FloatStruct(ctx context.Context, p1 Float) error
	FloatStructPointer(context.Context, *FloatPointer) error
	Bool(ctx context.Context, p1 bool) error
	BoolPointer(context.Context, *bool) error
	BoolArray(ctx context.Context, p1 []bool) error
	BoolArrayPointer(context.Context, []*bool) error
	BoolStruct(ctx context.Context, p1 Bool) error
	BoolStructPointer(context.Context, *BoolPointer) error
	Time(ctx context.Context, p1 time.Time) error
	TimePointer(context.Context, *time.Time) error
	TimeArray(ctx context.Context, p1 []time.Time) error
	TimeArrayPointer(context.Context, []*time.Time) error
	TimeShadow(ctx context.Context, p1 TimeShadow) error
	TimeShadowPointer(context.Context, *TimeShadow) error
	Decimal(ctx context.Context, p1 decimal.Decimal) error
	DecimalPinter(context.Context, *decimal.Decimal) error
	DecimalArray(ctx context.Context, p1 []decimal.Decimal) error
	DecimalArrayPointer(context.Context, []*decimal.Decimal) error
	DecimalShadow(ctx context.Context, p1 DecimalShadow) error
	DecimalShadowPointer(context.Context, *DecimalShadow) error
	Map(ctx context.Context, m1 map[string]string) error
	Type(ctx context.Context, p1 Type) error
	TypePointer(context.Context, *TypePointer) error
	TypeArray(ctx context.Context, p1 []Type) error
	TypeArrayPointer(context.Context, []*TypePointer) error
	TypeArrayArray(ctx context.Context, p1 [][]Type) error
	TypeArrayArrayPointer(context.Context, [][]*TypePointer) error
}

type ReplyAPI interface {
	//Int(ctx context.Context) (r1 int, err error)
	//IntPointer(ctx context.Context) (*int, error)
	//IntArray(ctx context.Context) (r1 []int, err error)
	//IntArrayPointer(ctx context.Context) ([]*int, error)
	//IntStruct(ctx context.Context) (r1 Int, err error)
	//IntStructPointer(ctx context.Context) (*IntPointer, error)
	//Uint(ctx context.Context) (r1 uint, err error)
	//UIntPointer(ctx context.Context) (*uint, error)
	//UintArray(ctx context.Context) (r1 []uint, err error)
	//UintArrayPointer(ctx context.Context) ([]*uint, error)
	//UintStruct(ctx context.Context) (r1 Uint, err error)
	//UintStructPointer(ctx context.Context) (*UintPointer, error)
	//String(ctx context.Context) (r1 string, err error)
	//StringPointer(ctx context.Context) (*string, error)
	//StringArray(ctx context.Context) (r1 []string, err error)
	//StringArrayPointer(ctx context.Context) ([]*string, error)
	//StringStruct(ctx context.Context) (r1 String, err error)
	//StringStructPointer(ctx context.Context) (*StringPointer, error)
	//Float(ctx context.Context) (r1 float32, err error)
	//FloatPointer(ctx context.Context) (*float32, error)
	//FloatArray(ctx context.Context) (r1 []float32, err error)
	//FloatArrayPointer(ctx context.Context) ([]*float32, error)
	//FloatStruct(ctx context.Context) (r1 Float, err error)
	//FloatStructPointer(ctx context.Context) (r1 *FloatPointer, err error)
	//Bool(ctx context.Context) (r1 bool, err error)
	//BoolPointer(ctx context.Context) (*bool, error)
	//BoolArray(ctx context.Context) (r1 []bool, err error)
	//BoolArrayPointer(ctx context.Context) ([]*bool, error)
	//BoolStruct(ctx context.Context) (r1 Bool, err error)
	//BoolStructPointer(ctx context.Context) (*BoolPointer, error)
	//Time(ctx context.Context) (r1 time.Time, err error)
	//TimePointer(ctx context.Context) (*time.Time, error)
	//TimeArray(ctx context.Context) (r1 []time.Time, err error)
	//TimeArrayPointer(ctx context.Context) ([]*time.Time, error)
	//TimeShadow(ctx context.Context) (r1 TimeShadow, err error)
	//TimeShadowPointer(ctx context.Context) (*TimeShadow, error)
	//Decimal(ctx context.Context) (r1 decimal.Decimal, err error)
	//DecimalPinter(ctx context.Context) (*decimal.Decimal, error)
	//DecimalArray(ctx context.Context) (r1 []decimal.Decimal, err error)
	//DecimalArrayPointer(ctx context.Context) ([]*decimal.Decimal, error)
	//DecimalShadow(ctx context.Context) (r1 DecimalShadow, err error)
	//DecimalShadowPointer(ctx context.Context) (*DecimalShadow, error)
	//Map(ctx context.Context) (m1 map[string]string, err error)
	//Type(ctx context.Context) (r1 Type, err error)
	//TypePointer(ctx context.Context) (*TypePointer, error)
	//TypeArray(ctx context.Context) (r1 []Type, err error)
	//TypeArrayPointer(context.Context) ([]*TypePointer, error)
	//TypeArrayArray(ctx context.Context) (r1 [][]Type, err error)
	//TypeArrayArrayPointer(context.Context) ([][]*TypePointer, error)
}

type Type struct {
	Int               Int
	Uint              Uint
	Float             Float
	Bool              Bool
	String            String
	Time              time.Time
	Decimal           decimal.Decimal
	IntArray          []Int
	UintArray         []Uint
	FloatArray        []Float
	BoolArray         []Bool
	StringArray       []String
	TimeArray         []time.Time
	DecimalArray      []decimal.Decimal
	IntArrayArray     [][]Int
	UintArrayArray    [][]Uint
	FloatArrayArray   [][]Float
	BoolArrayArray    [][]Bool
	StringArrayArray  [][]String
	TimeArrayArray    [][]time.Time
	DecimalArrayArray [][]decimal.Decimal
	Map               map[string]any
}

type TypePointer struct {
	Int               *IntPointer
	Uint              *UintPointer
	Float             *FloatPointer
	Bool              *BoolPointer
	String            *StringPointer
	Time              *time.Time
	Decimal           *decimal.Decimal
	IntArray          []*IntPointer
	UintArray         []*UintPointer
	FloatArray        []*FloatPointer
	BoolArray         []*BoolPointer
	StringArray       []*StringPointer
	TimeArray         []*time.Time
	DecimalArray      []*decimal.Decimal
	IntArrayArray     [][]*IntPointer
	UintArrayArray    [][]*UintPointer
	FloatArrayArray   [][]*FloatPointer
	BoolArrayArray    [][]*BoolPointer
	StringArrayArray  [][]*StringPointer
	TimeArrayArray    [][]*time.Time
	DecimalArrayArray [][]*decimal.Decimal
	Map               map[string]any
}

type NestType struct {
	Int
	Uint
	Float
	String
	Bool
}

type NestTypePointer struct {
	*IntPointer
	*UintPointer
	*FloatPointer
	*String
	*Bool
}

type NestNestType struct {
	NestType
}

type NestNestTypePointer struct {
	*NestTypePointer
}

type Int struct {
	Int             int
	Int8            int8
	Int16           int16
	Int32           int32
	Int64           int64
	IntArray        []int
	Int8Array       []int8
	Int16Array      []int16
	Int32Array      []int32
	Int64Array      []int64
	IntArrayArray   [][]int
	Int8ArrayArray  [][]int8
	Int16ArrayArray [][]int16
	Int32ArrayArray [][]int32
	Int64ArrayArray [][]int64
}

type IntPointer struct {
	Int             *int       `json:"int"`
	Int8            *int8      `json:"int8"`
	Int16           *int16     `json:"int16"`
	Int32           *int32     `json:"int32"`
	Int64           *int64     `json:"int64"`
	IntArray        []*int     `json:"intArray"`
	Int8Array       []*int8    `json:"int8Array"`
	Int16Array      []*int16   `json:"int16Array"`
	Int32Array      []*int32   `json:"int32Array"`
	Int64Array      []*int64   `json:"int64Array"`
	IntArrayArray   [][]*int   `json:"intArrayArray"`
	Int8ArrayArray  [][]*int8  `json:"int8ArrayArray"`
	Int16ArrayArray [][]*int16 `json:"int16ArrayArray"`
	Int32ArrayArray [][]*int32 `json:"int32ArrayArray"`
	Int64ArrayArray [][]*int64 `json:"int64ArrayArray"`
}

type Uint struct {
	Uint   uint
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64
}

type UintPointer struct {
	Uint   *uint   `json:"uint"`
	Uint8  *uint8  `json:"uint8"`
	Uint16 *uint16 `json:"uint16"`
	Uint32 *uint32 `json:"uint32"`
	Uint64 *uint64 `json:"uint64"`
}

type String struct {
	String string
}

type StringPointer struct {
	String *string `json:"string"`
}

type Bool struct {
	Bool bool
}

type BoolPointer struct {
	Bool *bool `json:"bool"`
}

type Float struct {
	Float32 float32
	Float64 float64
}

type FloatPointer struct {
	Float32 *float32 `json:"float32"`
	Float64 *float64 `json:"float64"`
}

type TimeShadow = time.Time

type DecimalShadow = decimal.Decimal
