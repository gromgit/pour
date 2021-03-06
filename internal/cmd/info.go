package cmd

import (
	"github.com/gromgit/pour/internal/console"
	"github.com/gromgit/pour/internal/formula"
	"log"
	"os"
	"strings"
	t "text/template"
)

type infoData struct {
	Me      formula.Formula
	All     formula.Formulas
	Deps    map[string]string
	Caveats map[string]string
}

var infoTemplates map[string]*t.Template

var funcMap = t.FuncMap{
	"status": func(status int) string {
		return formula.StatusMap[status]
	},
	"url": func(url string) string {
		return console.Underscore + url + console.Off
	},
}

func init() {
	infoTemplates = make(map[string]*t.Template)
	infoTemplates["main"] = t.Must(t.New("main").Funcs(funcMap).Parse(`{{.Me.Name}}: stable {{.Me.GetVersion}}
{{- if .Me.Pinned}} [pinned]{{end}}
{{.Me.Desc}}
{{url .Me.Homepage}}
{{if .Me.InstallDir -}} {{.Me.InstallDir}}{{status .Me.Status}}
  Poured from bottle on {{.Me.InstallTime}}
{{- else}}Not installed{{end}}
{{if .Me.Bottle.Stable.URL -}}From: {{url .Me.Bottle.Stable.URL}}{{- else}}No bottle found, cannot be installed{{end}}
{{if .Deps -}}
===> Dependencies
{{if .Deps.Required}}Required: {{.Deps.Required}}{{end -}}
{{if .Deps.Recommended}}Recommended: {{.Deps.Recommended}}{{end -}}
{{if .Deps.Optional}}Optional: {{.Deps.Optional}}{{end -}}
{{end}}
{{if .Caveats -}}
===> Caveats
{{if .Caveats.Specific}}
{{- .Caveats.Specific}}
{{- end}}
{{- if .Caveats.KegReason}}
{{.Caveats.Name}} is keg-only, which means it was not symlinked into {{.Caveats.Prefix}},
because {{.Caveats.KegReason}}.
{{end}}
{{- if .Caveats.Path}}
If you need to have {{.Caveats.Name}} first in your PATH run:
  echo 'export PATH="/usr/local/opt/{{.Caveats.Name}}/bin:$PATH"' >> ~/.bash_profile
{{end}}
{{- if .Caveats.Dev}}
For compilers to find {{.Caveats.Name}} you may need to set:
  export LDFLAGS="-L/usr/local/opt/{{.Caveats.Name}}/lib"
  export CPPFLAGS="-I/usr/local/opt/{{.Caveats.Name}}/include"
{{end}}
{{- if .Caveats.Pkgconfig}}
For pkg-config to find {{.Caveats.Name}} you may need to set:
  export PKG_CONFIG_PATH="/usr/local/opt/{{.Caveats.Name}}/lib/pkgconfig"
{{end}}
{{- end}}`))
}

func Info(allf formula.Formulas, args []string) error {
	tMain := infoTemplates["main"]
	for _, i := range args {
		if f := allf[i]; f.Name != "" {
			if err := tMain.Execute(
				os.Stdout,
				infoData{f,
					allf,
					f.GetDepReport(allf),
					f.GetCaveatReport()}); err != nil {
				log.Println("executing info template:", err)
				continue
			}
		}
	}
	return nil
}

func Leaves(allf formula.Formulas, args []string) error {
	allf.Filter(func(f formula.Formula) bool {
		return f.Leaf
	}).Ls()
	return nil
}

func Deps(allf formula.Formulas, args []string) error {
	installed := false
	common := false
	// Handle options
SearchOptions:
	for len(args) > 0 {
		switch {
		case strings.HasPrefix(args[0], "--inst"):
			installed = true
		case strings.HasPrefix(args[0], "--comm"):
			common = true
		default:
			break SearchOptions
		}
		args = args[1:]
	}
	// Filter for installed?
	if installed {
		var newArgs []string
		for _, n := range args {
			if allf[n].Installed() {
				newArgs = append(newArgs, n)
			}
		}
		args = newArgs
	}
	if len(args) > 0 {
		allf.Subset(allf.FindDeps(args, common)).Ls()
	}
	return nil
}

func Uses(allf formula.Formulas, args []string) error {
	installed := false
	recursive := false
	// Handle options
SearchOptions:
	for len(args) > 0 {
		switch {
		case strings.HasPrefix(args[0], "--inst"):
			installed = true
		case strings.HasPrefix(args[0], "--recur"):
			recursive = true
		default:
			break SearchOptions
		}
		args = args[1:]
	}
	if len(args) > 0 {
		result := allf.Subset(allf.FindUsers(args, recursive))
		// Filter for installed?
		if installed {
			result = result.Filter(func(f formula.Formula) bool {
				return f.Installed()
			})
		}
		result.Ls()
	}
	return nil
}
