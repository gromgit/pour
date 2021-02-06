package cmd

import (
	"errors"
	"github.com/gromgit/litebrew/internal/config"
	"os"
	fp "path/filepath"
	t "text/template"
)

var shellenvTemplates map[string]*t.Template

func init() {
	shellenvTemplates = make(map[string]*t.Template)
	shellenvTemplates["sh"] = t.Must(t.New("sh").Funcs(funcMap).Parse(`LITEBREW_PREFIX="{{.}}"; export LITEBREW_PREFIX;
LITEBREW_CELLAR="{{.}}/Cellar"; export LITEBREW_CELLAR;
LITEBREW_REPOSITORY="{{.}}/Lite"; export LITEBREW_REPOSITORY;
PATH="{{.}}/bin:{{.}}/sbin${PATH+:$PATH}"; export PATH;
MANPATH="{{.}}/share/man${MANPATH+:$MANPATH}:"; export MANPATH;
INFOPATH="{{.}}/share/info${INFOPATH+:$INFOPATH}"; export INFOPATH;
`))
	shellenvTemplates["bash"] = t.Must(t.New("bash").Funcs(funcMap).Parse(`export LITEBREW_PREFIX="{{.}}";
export LITEBREW_CELLAR="{{.}}/Cellar";
export LITEBREW_REPOSITORY="{{.}}/Lite";
export PATH="{{.}}/bin:{{.}}/sbin${PATH+:$PATH}";
export MANPATH="{{.}}/share/man${MANPATH+:$MANPATH}:";
export INFOPATH="{{.}}/share/info${INFOPATH+:$INFOPATH}";
`))
	shellenvTemplates["zsh"] = shellenvTemplates["bash"]
	shellenvTemplates["fish"] = t.Must(t.New("fish").Funcs(funcMap).Parse(`set -gx LITEBREW_PREFIX "{{.}}";
set -gx LITEBREW_CELLAR "{{.}}/Cellar";
set -gx LITEBREW_REPOSITORY "{{.}}/Homebrew";
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
