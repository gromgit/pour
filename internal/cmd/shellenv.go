package cmd

import (
	"errors"
	"github.com/gromgit/pour/internal/config"
	"os"
	fp "path/filepath"
	t "text/template"
)

var shellenvTemplates map[string]*t.Template

func init() {
	shellenvTemplates = make(map[string]*t.Template)
	shellenvTemplates["sh"] = t.Must(t.New("sh").Funcs(funcMap).Parse(`POUR_PREFIX="{{.}}"; export POUR_PREFIX;
POUR_CELLAR="{{.}}/Cellar"; export POUR_CELLAR;
PATH="{{.}}/bin:{{.}}/sbin${PATH+:$PATH}"; export PATH;
MANPATH="{{.}}/share/man${MANPATH+:$MANPATH}:"; export MANPATH;
INFOPATH="{{.}}/share/info${INFOPATH+:$INFOPATH}"; export INFOPATH;
`))
	shellenvTemplates["bash"] = t.Must(t.New("bash").Funcs(funcMap).Parse(`export POUR_PREFIX="{{.}}";
export POUR_CELLAR="{{.}}/Cellar";
export PATH="{{.}}/bin:{{.}}/sbin${PATH+:$PATH}";
export MANPATH="{{.}}/share/man${MANPATH+:$MANPATH}:";
export INFOPATH="{{.}}/share/info${INFOPATH+:$INFOPATH}";
`))
	shellenvTemplates["zsh"] = shellenvTemplates["bash"]
	shellenvTemplates["fish"] = t.Must(t.New("fish").Funcs(funcMap).Parse(`set -gx POUR_PREFIX "{{.}}";
set -gx POUR_CELLAR "{{.}}/Cellar";
set -g fish_user_paths "{{.}}/bin" "{{.}}/sbin" $fish_user_paths;
set -q MANPATH; or set MANPATH ''; set -gx MANPATH "{{.}}/share/man" $MANPATH;
set -q INFOPATH; or set INFOPATH ''; set -gx INFOPATH "{{.}}/share/info" $INFOPATH;
`))
}

func Shellenv(args []string) (err error) {
	shell := os.Getenv("SHELL")
	template := shellenvTemplates[fp.Base(shell)]
	if template == nil {
		err = errors.New("Unsupported shell " + shell)
	} else {
		err = template.Execute(os.Stdout, config.DEFAULT_PREFIX)
	}
	return
}
