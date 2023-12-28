package openapi

type API struct {
	Name    string    `json:"Name,omitempty"`
	Methods []*Method `json:"Methods,omitempty"`
	Defs    []*Def    `json:"Defs,omitempty"`
}

type Method struct {
	Name    string `json:"Name,omitempty"`
	Request *Def   `json:"Request,omitempty"`
	Replay  *Def   `json:"Replay,omitempty"`
	Comment string `json:"Comment,omitempty"`
}

type Def struct {
	Name         string `json:"Name,omitempty"`         // 类型声明时的名称
	Field        string `json:"Field,omitempty"`        // 作为 struct 字段时的名称
	Json         string `json:"Json,omitempty"`         // 作为 struct 字段时定义的 json
	Type         string `json:"Type"`                   // 类型，对标 golang 的类型
	Pointer      bool   `json:"Pointer,omitempty"`      // 是否为指针
	StructFields []*Def `json:"StructFields,omitempty"` // 存放 struct 字段
	ArrayFields  []*Def `json:"ArrayFields,omitempty"`  // 存放数组元素
	Elem         *Def   `json:"Elem,omitempty"`         // 存放数组引用值
	Use          *Def   `json:"Use,omitempty"`          // 存放类型引用
	Tags         string `json:"Tags,omitempty"`         // 存放 struct 字段时定义的标签
}
