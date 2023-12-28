package parser

type Method struct {
	Name    string
	Params  []*Param
	Results []*Param
}

func (m Method) ExportParams() []*Param {
	var n []*Param
	for _, v := range m.Params {
		if !v.Type.IsContext() {
			n = append(n, v)
		}
	}
	return n
}

func (m Method) ExportResults() []*Param {
	var n []*Param
	for _, v := range m.Results {
		if !v.Type.IsError() {
			n = append(n, v)
		}
	}
	return n
}
