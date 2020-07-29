package dep

import (
	"fmt"
	"path/filepath"

	gosrc "github.com/Morganamilo/go-srcinfo"
)

type LocalPkg struct {
	*gosrc.Srcinfo
}

func NewLocalPackage(path string) (*LocalPkg, error) {
	pkgbuild, err := gosrc.ParseFile(filepath.Join(path, ".SRCINFO"))
	if err != nil {
		return nil, err
	}
	fmt.Println(pkgbuild)
	return &LocalPkg{pkgbuild}, nil
}

func (p LocalPkg) Name() string {
	return p.Srcinfo.Packages[0].Pkgname
}

func (p LocalPkg) Version() string {
	return p.Srcinfo.Pkgver
}

func (p LocalPkg) PackageBase() string {
	return p.Srcinfo.Pkgbase
}

func (p LocalPkg) Provides() []string {
	return archStringToString(p.Srcinfo.Provides)
}

func (p LocalPkg) Depends() []string {
	return archStringToString(p.Srcinfo.Depends)
}

func (p LocalPkg) MakeDepends() []string {
	return archStringToString(p.Srcinfo.MakeDepends)
}

func (p LocalPkg) CheckDepends() []string {
	return archStringToString(p.Srcinfo.CheckDepends)
}

func archStringToString(as []gosrc.ArchString) []string {
	var provides []string
	for _, p := range as {
		provides = append(provides, p.Value)
	}
	return provides
}
