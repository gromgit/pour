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

type InstallMap map[string]Action

func resolveDependencies(allf *formula.Formulas, name string, depMap InstallMap, leaf bool) {
	// Only "local" bottles need apply
	if strings.Contains(name, "/") {
		depMap[name] = Action{ERROR, "Cannot install foreign bottle"}
		return
	}
	f := (*allf)[name]
	if f.Name == "" {
		depMap[name] = Action{ERROR, "Bottle not found"}
		return
	}
	s := NOTHING
	switch f.Status {
	case formula.MISSING:
		s = INSTALL
	case formula.OUTDATED:
		if f.Pinned {
			depMap[name] = Action{ERROR, "Cannot update pinned bottle"}
			return
		} else {
			s = UPGRADE
		}
	}
	if leaf {
		s |= LEAF
	}
	if depMap[name].Code == NOTHING {
		depMap[name] = Action{s, ""}
	}
	for _, dep := range f.Dependencies {
		resolveDependencies(allf, dep, depMap, false)
	}
	depMap[name] = Action{s, ""}
	return
}

type ActionPlan struct {
	Name string
	Leaf bool
}

func Install(allf *formula.Formulas, args []string) (err error) {
	instMap := make(InstallMap)
	for _, name := range args {
		resolveDependencies(allf, name, instMap, true)
	}
	// Do all the installations
	log.Logf("instMap: %+v\n", instMap)
	var errors []string
	var actions []ActionPlan
	for name, act := range instMap {
		switch act.Code & ACT_MASK {
		case ERROR:
			errors = append(errors, name+": "+act.Message)
		case INSTALL, UPGRADE:
			actions = append(actions, ActionPlan{name, act.Code&LEAF > 0})
		}
	}
	if len(errors) == 0 {
		for _, a := range actions {
			if err = bottle.Install((*allf)[a.Name], a.Leaf); err != nil {
				errors = append(errors, a.Name+": "+err.Error())
			}
		}
	}
	if len(errors) > 0 {
		fmt.Println("===> ERRORS\n" + strings.Join(errors, "\n"))
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
