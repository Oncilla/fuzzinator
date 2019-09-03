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

package execute

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/otiai10/copy"
	"golang.org/x/xerrors"
	"gopkg.in/src-d/go-git.v4"

	"github.com/Oncilla/fuzzinator/conf"
)

// 1. Setup corpus
//  a. Create tmp dir with current commit + package as identifier
//  b. Copy corpus to tmp dir
// 2. Build binary
// 3. Run binary
// 4. Wait for SIGINT, commit crashers

// commit, err := CommitHash()
// if err != nil {
// 	return xerrors.Errorf("unable to get commit hash: %w", err)
// }

// SetupTempDir sets up the temporary working directory and returns the path.
func SetupTempDir(targetName, commit string) (string, error) {
	tmpDir := TempCorpusName(targetName, commit)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return "", xerrors.Errorf("unable to create temporary corpus: %w", err)
	}
	return tmpDir, nil
}

// SetupCorpus sets up the temporary working directory with the configured corpus.
func SetupCorpus(corpus, tmpDir string) error {
	if err := copy.Copy(corpus, filepath.Join(tmpDir, "corpus")); err != nil {
		return xerrors.Errorf("unable to copy corpus: %w", err)
	}
	return nil
}

// BuildBinary builds the fuzzing binary and returns the path.
func BuildBinary(target conf.Target, tmpDir string) (string, error) {
	output := filepath.Join(tmpDir, "fuzz.zip")
	cmd := exec.Command("go-fuzz-build", "-o", output, "-tags", target.Harness.BuildTags,
		"-func", target.Harness.Function, target.Harness.Package)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", xerrors.Errorf("unable to build fuzzing binary: %x", err)
	}
	return output, nil
}

// go-fuzz -bin=./bin/lib_hpkt-fuzz.zip -workdir=workdir/lib_hpkt
// go-fuzz-build -o bin/lib_hpkt-fuzz.zip ./lib_hpkt

// TempCorpusName returns the temporary corpus name for a given target and commit.
func TempCorpusName(target string, commit string) string {
	return filepath.Join(os.TempDir(), "fuzzinator", fmt.Sprintf("%s_%s", target, commit))
}

// CommitHash gets the commit hash of the local git repository.
func CommitHash() (string, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return "", xerrors.Errorf("unable to open git repository: %w", err)
	}
	ref, err := r.Head()
	if err != nil {
		return "", xerrors.Errorf("unable to determine head: %w", err)
	}
	return ref.Hash().String(), nil
}
