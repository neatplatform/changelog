package main

import (
	"context"
	"fmt"
	"os"

	"github.com/neatplatform/craft/flagit"
	"github.com/neatplatform/craft/ui"

	"github.com/neatplatform/changelog/generate"
	"github.com/neatplatform/changelog/internal/git"
	"github.com/neatplatform/changelog/metadata"
	"github.com/neatplatform/changelog/spec"
)

func main() {
	// Verbosity level will be updated once it is known.
	u := ui.New(ui.None)

	/* -------------------- READING SPEC -------------------- */

	s, err := spec.Default().FromFile()
	if err != nil {
		u.Errorf(ui.Red, "%s", err)
		os.Exit(1)
	}

	if err := flagit.Parse(&s, false); err != nil {
		u.Errorf(ui.Red, "%s", err)
		os.Exit(1)
	}

	// Update the verbosity level.
	if s.General.Verbose {
		u.SetLevel(ui.Debug)
	} else if !s.General.Print {
		u.SetLevel(ui.Info)
	}

	u.Debugf(ui.Cyan, "%s", s)

	/* -------------------- RUNNING COMMANDS -------------------- */

	switch {
	case s.Help:
		if err := s.PrintHelp(); err != nil {
			u.Errorf(ui.Red, "%s", err)
			os.Exit(1)
		}

	case s.Version:
		fmt.Println(metadata.String())

	default:
		gitRepo, err := git.NewRepo(u, ".")
		if err != nil {
			u.Errorf(ui.Red, "%s", err)
			os.Exit(1)
		}

		domain, path, err := gitRepo.GetRemote()
		if err != nil {
			u.Errorf(ui.Red, "%s", err)
			os.Exit(1)
		}

		s = s.WithRepo(domain, path)

		g, err := generate.New(s, u)
		if err != nil {
			u.Errorf(ui.Red, "%s", err)
			os.Exit(1)
		}

		ctx := context.Background()

		if _, err := g.Generate(ctx, s); err != nil {
			u.Errorf(ui.Red, "%s", err)
			os.Exit(1)
		}
	}
}
