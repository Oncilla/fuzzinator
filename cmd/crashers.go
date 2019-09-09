// MIT License
//
// Copyright (c) 2019 Oncilla
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"github.com/oncilla/fuzzinator/lib"
)

const msgTmplFmt = `
Add crashers:
  - target: "%s"
  - pkg:    "%s"
  - entry:  "%s"
  - commit: "%s"
`

var crashersCmd = &cobra.Command{
	Use:   "crashers",
	Short: "copy the crashers to the corpus and commit them",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target, commit, err := targetAndCommit(confFile, args[0])
		if err != nil {
			return err
		}
		crashers := crashersOut(target.Corpus, commit, target.Crashers)
		if err := copyCrashers(target.Name, crashers, commit, terminate); err != nil {
			return err
		}
		added, err := lib.AddCrashers(crashers, target.Name, commit)
		if err != nil {
			return err
		}
		if !added {
			log.Println("No new crashers added")
			return nil
		}
		msg := fmt.Sprintf(msgTmplFmt, target.Name, target.Harness.Package,
			target.Harness.Function, commit)
		log.Println("Please commit added crashers:")
		fmt.Printf("\ngit commit -m '%s'\n", msg)
		return nil
	},
}

func copyCrashers(name, crashers, commit string, stop <-chan struct{}) error {
	if err := os.MkdirAll(crashers, 0755); err != nil {
		return xerrors.Errorf("unable to create crashers dir: %w", err)
	}
	log.Printf("Copying crashers to %q", crashers)
	if err := lib.CopyCrashers(lib.TempWorkdir(name, commit), crashers); err != nil {
		return err
	}
	return nil
}

func crashersOut(corpus, commit, crashers string) string {
	if crashers != "" {
		return filepath.Join(crashers, commit)
	}
	return filepath.Join(filepath.Dir(corpus), "crashers", commit)
}
