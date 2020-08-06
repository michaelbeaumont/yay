package dep

type Pkg interface {
	Name() string
	Version() string
	PackageBase() string
	Provides() []string
	Depends() []string
	MakeDepends() []string
	CheckDepends() []string
}

type Base interface {
	Pkgbase() string
	Version() string
	String() string
	Pkgs() []Pkg
}
