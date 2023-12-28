package parser

func Parse(mod *Mod, dir string) (pkg *Package, err error) {
	pkg = &Package{}
	err = pkg.Parse(mod, dir)
	if err != nil {
		return
	}
	return
}
