package cmd

import (
	"fmt"
	"github.com/gromgit/litebrew/internal/config"
)

func Shellenv(args []string) error {
	fmt.Printf(`export LITEBREW_PREFIX="%s";
export LITEBREW_CELLAR="${LITEBREW_PREFIX}/Cellar";
export LITEBREW_REPOSITORY="${LITEBREW_PREFIX}/Lite";
export PATH="${LITEBREW_PREFIX}/bin:${LITEBREW_PREFIX}/sbin${PATH+:$PATH}";
export MANPATH="${LITEBREW_PREFIX}/share/man${MANPATH+:$MANPATH}:";
export INFOPATH="${LITEBREW_PREFIX}/share/info${INFOPATH+:$INFOPATH}";
`, config.DEFAULT_PREFIX)
	return nil
}
