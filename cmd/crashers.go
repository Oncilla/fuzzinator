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
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"github.com/oncilla/fuzzinator/conf"
	"github.com/oncilla/fuzzinator/lib"
)

var crashersCmd = &cobra.Command{
	Use:   "crashers",
	Short: "copy the crashers to the corpus",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target, commit, err := targetAndCommit(confFile, args[0])
		if err != nil {
			return err
		}
		return crashersCrashers(target, commit, terminate)
	},
}

func crashersCrashers(target conf.Target, commit string, stop <-chan struct{}) error {
	crashers := filepath.Join(target.Crashers, commit)
	if target.Crashers == "" {
		crashers = filepath.Join(filepath.Dir(target.Corpus), "crashers", commit)
	}
	if err := os.MkdirAll(crashers, 0755); err != nil {
		return xerrors.Errorf("unable to create crashers dir: %w", err)
	}
	log.Println("Copying crashers")
	if err := lib.CopyCrashers(lib.TempWorkdir(target.Name, commit), crashers); err != nil {
		return err
	}
	return nil
}
