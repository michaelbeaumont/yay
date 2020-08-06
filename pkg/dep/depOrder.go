package dep

import (
	"fmt"

	alpm "github.com/Jguer/go-alpm"

	"github.com/Jguer/yay/v10/pkg/stringset"
	"github.com/Jguer/yay/v10/pkg/text"
)

type Order struct {
	Aur     []AURBase
	Local   []LocalBase
	Repo    []*alpm.Package
	Runtime stringset.StringSet
}

func makeOrder() *Order {
	return &Order{
		make([]AURBase, 0),
		make([]LocalBase, 0),
		make([]*alpm.Package, 0),
		make(stringset.StringSet),
	}
}

func GetOrder(dp *Pool) *Order {
	do := makeOrder()

	for _, target := range dp.Targets {
		dep := target.DepString()
		if aurPkg := dp.Aur[dep]; aurPkg != nil && pkgSatisfies(aurPkg.Name, aurPkg.Version, dep) {
			do.orderPkgAur(RPCPkg{aurPkg}, dp, true)
		}

		if localPkg := dp.findSatisfierLocal(dep); localPkg != nil {
			do.orderPkgLocal(localPkg, dp, true)
		}

		if aurPkg := dp.findSatisfierAur(dep); aurPkg != nil {
			do.orderPkgAur(aurPkg, dp, true)
		}

		if repoPkg := dp.findSatisfierRepo(dep); repoPkg != nil {
			do.orderPkgRepo(repoPkg, dp, true)
		}
	}

	return do
}

func (do *Order) orderPkgAur(pkg Pkg, dp *Pool, runtime bool) {
	if runtime {
		do.Runtime.Set(pkg.Name())
	}
	delete(dp.Aur, pkg.Name())

	for i, deps := range [3][]string{pkg.Depends(), pkg.MakeDepends(), pkg.CheckDepends()} {
		for _, dep := range deps {
			aurPkg := dp.findSatisfierAur(dep)
			if aurPkg != nil {
				do.orderPkgAur(aurPkg, dp, runtime && i == 0)
			}

			repoPkg := dp.findSatisfierRepo(dep)
			if repoPkg != nil {
				do.orderPkgRepo(repoPkg, dp, runtime && i == 0)
			}
		}
	}

	for i, base := range do.Aur {
		if base.Pkgbase() == pkg.PackageBase() {
			do.Aur[i] = append(base, pkg)
			return
		}
	}

	do.Aur = append(do.Aur, AURBase{pkg})
}

func (do *Order) orderPkgLocal(pkg *LocalPkg, dp *Pool, runtime bool) {
	if runtime {
		do.Runtime.Set(pkg.Name())
	}
	delete(dp.Local, pkg.Name())

	for i, deps := range [3][]string{pkg.Depends(), pkg.MakeDepends(), pkg.CheckDepends()} {
		for _, dep := range deps {
			aurPkg := dp.findSatisfierAur(dep)
			if aurPkg != nil {
				do.orderPkgAur(aurPkg, dp, runtime && i == 0)
			}

			repoPkg := dp.findSatisfierRepo(dep)
			if repoPkg != nil {
				do.orderPkgRepo(repoPkg, dp, runtime && i == 0)
			}
		}
	}

	for i, base := range do.Local {
		if base.Pkgbase == pkg.PackageBase() {
			do.Local[i] = base
			return
		}
	}

	do.Local = append(do.Local, pkg.AsBase())
}

func (do *Order) orderPkgRepo(pkg *alpm.Package, dp *Pool, runtime bool) {
	if runtime {
		do.Runtime.Set(pkg.Name())
	}
	delete(dp.Repo, pkg.Name())

	_ = pkg.Depends().ForEach(func(dep alpm.Depend) (err error) {
		repoPkg := dp.findSatisfierRepo(dep.String())
		if repoPkg != nil {
			do.orderPkgRepo(repoPkg, dp, runtime)
		}

		return nil
	})

	do.Repo = append(do.Repo, pkg)
}

func (do *Order) HasMake() bool {
	lenPkgs := 0
	for _, base := range do.Aur {
		lenPkgs += len(base)
	}
	// TODO use base
	lenPkgs = lenPkgs + len(do.Local)

	return len(do.Runtime) != lenPkgs+len(do.Repo)
}

func (do *Order) GetMake() []string {
	makeOnly := []string{}

	for _, base := range do.Aur {
		for _, pkg := range base {
			if !do.Runtime.Get(pkg.Name()) {
				makeOnly = append(makeOnly, pkg.Name())
			}
		}
	}

	for _, base := range do.Local {
		for _, pkg := range base.Pkgs() {
			if !do.Runtime.Get(pkg.Name()) {
				makeOnly = append(makeOnly, pkg.Name())
			}
		}
	}

	for _, pkg := range do.Repo {
		if !do.Runtime.Get(pkg.Name()) {
			makeOnly = append(makeOnly, pkg.Name())
		}
	}

	return makeOnly
}

// Print prints repository packages to be downloaded
func (do *Order) Print() {
	repo := ""
	repoMake := ""
	aur := ""
	aurMake := ""

	repoLen := 0
	repoMakeLen := 0
	aurLen := 0
	aurMakeLen := 0

	for _, pkg := range do.Repo {
		pkgStr := fmt.Sprintf("  %s-%s", pkg.Name(), pkg.Version())
		if do.Runtime.Get(pkg.Name()) {
			repo += pkgStr
			repoLen++
		} else {
			repoMake += pkgStr
			repoMakeLen++
		}
	}

	bases := []Base{}
	for _, p := range do.Aur {
		bases = append(bases, p)
	}
	for _, p := range do.Local {
		bases = append(bases, p)
	}
	for _, base := range bases {
		pkg := base.Pkgbase()
		pkgStr := "  " + pkg + "-" + base.Pkgs()[0].Version()
		pkgStrMake := pkgStr

		push := false
		pushMake := false

		switch {
		case len(base.Pkgs()) > 1, pkg != base.Pkgs()[0].Name():
			pkgStr += " ("
			pkgStrMake += " ("

			for _, split := range base.Pkgs() {
				if do.Runtime.Get(split.Name()) {
					pkgStr += split.Name() + " "
					aurLen++
					push = true
				} else {
					pkgStrMake += split.Name() + " "
					aurMakeLen++
					pushMake = true
				}
			}

			pkgStr = pkgStr[:len(pkgStr)-1] + ")"
			pkgStrMake = pkgStrMake[:len(pkgStrMake)-1] + ")"
		case do.Runtime.Get(base.Pkgs()[0].Name()):
			aurLen++
			push = true
		default:
			aurMakeLen++
			pushMake = true
		}

		if push {
			aur += pkgStr
		}
		if pushMake {
			aurMake += pkgStrMake
		}
	}

	printDownloads("Repo", repoLen, repo)
	printDownloads("Repo Make", repoMakeLen, repoMake)
	printDownloads("Aur", aurLen, aur)
	printDownloads("Aur Make", aurMakeLen, aurMake)
}

func printDownloads(repoName string, length int, packages string) {
	if length < 1 {
		return
	}

	repoInfo := fmt.Sprintf(text.Bold(text.Blue("[%s:%d]")), repoName, length)
	fmt.Println(repoInfo + text.Cyan(packages))
}