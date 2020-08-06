package dep

import (
	"fmt"
	"path/filepath"

	gosrc "github.com/Morganamilo/go-srcinfo"
)

type LocalBase struct {
	srcinfo *gosrc.Srcinfo
}

func NewLocalBase(path string) (LocalBase, error) {
	pkgbuild, err := gosrc.ParseFile(filepath.Join(path, ".SRCINFO"))
	if err != nil {
		return LocalBase{}, err
	}
	fmt.Println(pkgbuild)
	return LocalBase{pkgbuild}, nil
}

func (b LocalBase) Pkgbase() string {
	return b.srcinfo.Pkgbase
}
func (b LocalBase) Version() string {
	return b.srcinfo.Pkgver
}
func (b LocalBase) String() string {
	return b.srcinfo.String()
}

func (b LocalBase) Pkgs() []Pkg {
	bs := make([]Pkg, 0, len(b.srcinfo.Packages))
	for i := range b.srcinfo.Packages {
		bs = append(bs, LocalPkg{
			srcinfo: b.srcinfo,
			index:   i,
		})
	}
	return bs
}

type LocalPkg struct {
	srcinfo *gosrc.Srcinfo
	index   int
}

func (p LocalPkg) AsBase() LocalBase {
	return LocalBase{
		p.srcinfo,
	}
}

func (p LocalPkg) Name() string {
	return p.srcinfo.Packages[p.index].Pkgname
}

func (p LocalPkg) Version() string {
	return p.srcinfo.Pkgver
}

func (p LocalPkg) PackageBase() string {
	return p.srcinfo.Pkgbase
}

func (p LocalPkg) Provides() []string {
	return archStringToString(p.srcinfo.Packages[p.index].Provides)
}

func (p LocalPkg) Depends() []string {
	return archStringToString(p.srcinfo.Packages[p.index].Depends)
}

func (p LocalPkg) MakeDepends() []string {
	return archStringToString(p.srcinfo.MakeDepends)
}

func (p LocalPkg) CheckDepends() []string {
	return archStringToString(p.srcinfo.CheckDepends)
}

func archStringToString(as []gosrc.ArchString) []string {
	var provides []string
	for _, p := range as {
		provides = append(provides, p.Value)
	}
	return provides
}
