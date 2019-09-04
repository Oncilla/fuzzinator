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
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"github.com/Oncilla/fuzzinator/conf"
	"github.com/Oncilla/fuzzinator/lib"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup the temporary workdir and build the fuzzing binary",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target, commit, err := targetAndCommit(confFile, args[0])
		if err != nil {
			return err
		}
		return setup(target, commit, terminate)
	},
}

func setup(target conf.Target, commit string, stop <-chan struct{}) error {
	tmpDir, err := lib.SetupTempWorkdir(target.Name, commit)
	if err != nil {
		xerrors.Errorf("unable to setup temp dir: %w", err)
	}
	if err := lib.SetupCorpus(target.Corpus, tmpDir); err != nil {
		xerrors.Errorf("unable to setup corpus: %w", err)
	}
	if _, err := lib.BuildBinary(target, tmpDir, stop); err != nil {
		xerrors.Errorf("unable to build fuzzing binary: %w", err)
	}
	return nil
}
