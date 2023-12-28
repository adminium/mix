package parser

type Import struct {
	Name     string
	location string
	*Package
}
