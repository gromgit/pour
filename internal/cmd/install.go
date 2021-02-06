package cmd

import (
	"fmt"
	"github.com/gromgit/pour/internal/bottle"
	cfg "github.com/gromgit/pour/internal/config"
	"github.com/gromgit/pour/internal/formula"
	"github.com/gromgit/pour/internal/log"
	"path/filepath"
	"strings"
)

const (
	NOTHING = iota
	UPGRADE
	INSTALL
	// Leaf bit
	LEAF     = 1 << iota
	ACT_MASK = LEAF - 1
)

type InstallMap map[string]int

func resolveDependencies(allf *formula.Formulas, name string, depMap InstallMap, leaf bool) error {
	// Only "local" bottles need apply
	if strings.Contains(name, "/") {
		return fmt.Errorf("Unable to install foreign bottle '%s'", name)
	}
	f := (*allf)[name]
	if f.Name == "" {
		return fmt.Errorf("No such bottle '%s'", name)
	}
	s := NOTHING
	switch f.Status {
	case formula.MISSING:
		s = INSTALL
	case formula.OUTDATED:
		if f.Pinned {
			return fmt.Errorf("Can't update pinned bottle '%s'", name)
		} else {
			s = UPGRADE
		}
	}
	if leaf {
		s |= LEAF
	}
	if depMap[name] == NOTHING {
		depMap[name] = s
	}
	for _, dep := range f.Dependencies {
		if err := resolveDependencies(allf, dep, depMap, false); err != nil {
			return err
		}
	}
	return nil
}

func Install(allf *formula.Formulas, args []string) (err error) {
	instMap := make(InstallMap)
	for _, name := range args {
		if err = resolveDependencies(allf, name, instMap, true); err != nil {
			return
		}
	}
	// Do all the installations
	log.Logf("instMap: %+v\n", instMap)
	for name, act := range instMap {
		f := (*allf)[name]
		if act&ACT_MASK != NOTHING {
			if err = bottle.Install(f, (act&LEAF > 0)); err != nil {
				return
			}
		}
	}
	return
}

func Link(allf *formula.Formulas, args []string) (err error) {
	for _, name := range args {
		f := (*allf)[name]
		if f.Installed() && f.InstallDir != "" {
			rel, err := filepath.Rel(cfg.CELLAR, f.InstallDir)
			if err != nil {
				return err
			}
			err = bottle.Link(rel)
			if err != nil {
				return err
			}
		}
	}
	return
}
