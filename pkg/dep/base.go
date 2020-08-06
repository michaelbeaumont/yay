package dep

import rpc "github.com/mikkeloscar/aur"

type RPCPkg struct {
	*rpc.Pkg
}

func (p RPCPkg) Name() string {
	return p.Pkg.Name
}

func (p RPCPkg) Version() string {
	return p.Pkg.Version
}

func (p RPCPkg) PackageBase() string {
	return p.Pkg.PackageBase
}

func (p RPCPkg) Provides() []string {
	return p.Pkg.Provides
}

func (p RPCPkg) Depends() []string {
	return p.Pkg.Depends
}

func (p RPCPkg) MakeDepends() []string {
	return p.Pkg.MakeDepends
}

func (p RPCPkg) CheckDepends() []string {
	return p.Pkg.CheckDepends
}

// Base is an AUR base package
type AURBase []Pkg

// Pkgbase returns the first base package.
func (b AURBase) Pkgbase() string {
	return b[0].PackageBase()
}

// Version returns the first base package version.
func (b AURBase) Version() string {
	return b[0].Version()
}

// Packages foo and bar from a pkgbase named base would print like so:
// base (foo bar)
func (b AURBase) String() string {
	pkg := b[0]
	str := pkg.PackageBase()
	if len(b) > 1 || pkg.PackageBase() != pkg.Name() {
		str2 := " ("
		for _, split := range b {
			str2 += split.Name() + " "
		}
		str2 = str2[:len(str2)-1] + ")"

		str += str2
	}

	return str
}

func (b AURBase) Pkgs() []Pkg {
	return b
}

func GetBases(pkgs []*rpc.Pkg) []Base {
	basesMap := make(map[string]AURBase)
	for _, pkg := range pkgs {
		basesMap[pkg.PackageBase] = append(basesMap[pkg.PackageBase], RPCPkg{pkg})
	}

	bases := make([]Base, 0, len(basesMap))
	for _, base := range basesMap {
		bases = append(bases, base)
	}

	return bases
}
