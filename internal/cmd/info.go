package cmd

import (
	"github.com/gromgit/litebrew/internal/console"
	"github.com/gromgit/litebrew/internal/formula"
	"log"
	"os"
	t "text/template"
)

type templateMap map[string]*t.Template
type templateData struct {
	Me      formula.Formula
	All     formula.Formulas
	Deps    map[string]string
	Caveats map[string]string
}

var templates templateMap

var funcMap = t.FuncMap{
	"status": func(status int) string {
		return formula.StatusMap[status]
	},
	"url": func(url string) string {
		return console.Underscore + url + console.Off
	},
}

func init() {
	templates = make(templateMap)
	templates["main"] = t.Must(t.New("main").Funcs(funcMap).Parse(`{{.Me.Name}}: stable {{.Me.GetVersion}}
{{- if .Me.Pinned}} [pinned]{{end}}
{{.Me.Desc}}
{{url .Me.Homepage}}
{{if .Me.InstallDir -}} {{.Me.InstallDir}}{{status .Me.Status}}
  Poured from bottle on {{.Me.InstallTime}}
{{- else}}Not installed{{end}}
From: {{url .Me.Bottle.Stable.URL}}
{{- if .Deps}}
===> Dependencies
{{if .Deps.Required}}Required: {{.Deps.Required}}{{end -}}
{{if .Deps.Recommended}}Recommended: {{.Deps.Recommended}}{{end -}}
{{if .Deps.Optional}}Optional: {{.Deps.Optional}}{{end -}}
{{end}}
{{- if .Caveats}}
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

func Info(formulas formula.Formulas, args []string) {
	tMain := templates["main"]
	for _, i := range args {
		if f := formulas[i]; f.Name != "" {
			err := tMain.Execute(os.Stdout, templateData{f, formulas, f.GetDepReport(formulas), f.GetCaveatReport()})
			if err != nil {
				log.Println("executing main template:", err)
				continue
			}
		}
	}
}
