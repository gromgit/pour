package cmd

import (
	"github.com/gromgit/litebrew/internal/console"
	"github.com/gromgit/litebrew/internal/formula"
	"log"
	"os"
	t "text/template"
)

type templateMap map[string]*t.Template

var templates templateMap

var funcMap = t.FuncMap{
	"status": func(status int) string { return formula.StatusMap[status] },
	"url":    func(url string) string { return console.Underscore + url + console.Off },
}

func init() {
	templates = make(templateMap)
	templates["main"] = t.Must(t.New("main").Funcs(funcMap).Parse(`{{.Name}}: stable {{.GetVersion}}{{if .Pinned}} [pinned]{{end}}
{{.Desc}}
{{url .Homepage}}
{{.InstallDir}}{{status .Status}}
From: {{url .Bottle.Stable.URL}}
`))
}

func Info(formulas formula.Formulas, args []string) {
	t := templates["main"]
	for _, i := range args {
		if f := formulas[i]; f.Name != "" {
			err := t.Execute(os.Stdout, f)
			if err != nil {
				log.Println("executing template:", err)
			}
		}
	}
}
