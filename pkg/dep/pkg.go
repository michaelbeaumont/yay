package dep

import rpc "github.com/mikkeloscar/aur"

type Pkg interface {
	Name() string
	Version() string
	PackageBase() string
	Provides() []string
	Depends() []string
	MakeDepends() []string
	CheckDepends() []string
}

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
