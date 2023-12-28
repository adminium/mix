package parser

type Param struct {
	Names []string
	Type  *Type
	Def   *Def `json:"Def,omitempty"`
}
