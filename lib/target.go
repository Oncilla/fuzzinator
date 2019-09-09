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

package lib

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/otiai10/copy"
	"golang.org/x/tools/go/packages"
	"golang.org/x/xerrors"
	"gopkg.in/src-d/go-git.v4"

	"github.com/oncilla/fuzzinator/conf"
)

// SetupTempWorkdir sets up the temporary working directory and returns the path.
func SetupTempWorkdir(targetName, commit string) (string, error) {
	workdir := TempWorkdir(targetName, commit)
	if err := os.MkdirAll(workdir, 0755); err != nil {
		return "", xerrors.Errorf("unable to create temporary corpus: %w", err)
	}
	return workdir, nil
}

// TempWorkdir returns the temporary workdir path for a given target and commit.
func TempWorkdir(targetName string, commit string) string {
	return filepath.Join(os.TempDir(), "fuzzinator", fmt.Sprintf("%s_%s", targetName, commit))
}

// SetupCorpus sets up the temporary working directory with the configured corpus.
func SetupCorpus(corpus, workdir string) error {
	if err := copy.Copy(corpus, filepath.Join(workdir, "corpus")); err != nil {
		return xerrors.Errorf("unable to copy corpus: %w", err)
	}
	return nil
}

// CopyCrashers copies the crashers from the workdir to the target directory.
func CopyCrashers(workdir, target string) error {
	if err := copy.Copy(filepath.Join(workdir, "crashers"), target); err != nil {
		return xerrors.Errorf("unable to copy crashers: %w", err)
	}
	return nil
}

// BuildBinary builds the fuzzing binary and returns the path.
func BuildBinary(target conf.Target, workdir string, stop <-chan struct{}) (string, error) {
	output := BinaryPath(workdir)
	cmd := exec.Command("go-fuzz-build", "-o", output, "-tags", target.Harness.BuildTags,
		"-func", target.Harness.Function, target.Harness.Package)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return "", xerrors.Errorf("unable to start building fuzzing binary: %w", err)
	}
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-stop:
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			return "", xerrors.Errorf("unable to terminate building fuzzing binary: %w", err)
		}
		return "", xerrors.Errorf("abort building due to SIGTERM")
	case err := <-done:
		if err != nil {
			return "", xerrors.Errorf("error while bulding fuzzing binary: %w", err)
		}
		return output, nil
	}
}

// BinaryPath returns the file path to the fuzzing binary based on the temporary directory.
func BinaryPath(workdir string) string {
	return filepath.Join(workdir, "fuzz.zip")
}

// RunBinary runs the fuzzing binary until the stop channel is closed.
func RunBinary(fuzzBin string, workdir string, stop <-chan struct{}) error {
	cmd := exec.Command("go-fuzz", "-bin", fuzzBin, "-workdir", workdir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return xerrors.Errorf("unable to build fuzzing binary: %w", err)
	}
	<-stop
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return xerrors.Errorf("unable to terminate fuzzing: %w", err)
	}
	return nil
}

// PkgDir returns the absolute path to a go package.
func PkgDir(pkg string) (string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedFiles,
	}
	respkgs, err := packages.Load(cfg, pkg)
	if err != nil {
		return "", xerrors.Errorf("unable to resolve package: %w", err)
	}
	if len(respkgs) != 1 {
		paths := make([]string, len(respkgs))
		for i, p := range respkgs {
			paths[i] = p.PkgPath
		}
		return "", xerrors.Errorf("cannot build multiple packages, but %q resolved to: %v",
			pkg, strings.Join(paths, ", "))
	}
	info := respkgs[0]
	if len(info.GoFiles) == 0 {
		return "", xerrors.Errorf("no go file for %q", pkg)
	}
	return filepath.Dir(info.GoFiles[0]), nil
}

// CommitHash gets the commit hash of the local git repository.
func CommitHash(dir string) (string, error) {
	r, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return "", xerrors.Errorf("unable to open git repository at %q: %w", dir, err)
	}
	ref, err := r.Head()
	if err != nil {
		return "", xerrors.Errorf("unable to determine head: %w", err)
	}
	return ref.Hash().String(), nil
}

// AddCrashers adds the crashers to the git repository.
func AddCrashers(crashers, name, commit string) (bool, error) {
	r, err := git.PlainOpenWithOptions(crashers, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return false, xerrors.Errorf("unable to open git repository at %q: %w", crashers, err)
	}
	w, err := r.Worktree()
	if err != nil {
		return false, xerrors.Errorf("unable to get worktree: %w", err)
	}
	s, err := w.Status()
	if err != nil {
		return false, xerrors.Errorf("cannot determine status: %w", err)
	}
	if err := isClean(s); err != nil {
		return false, xerrors.Errorf("can only auto-commit on clean worktree: %s", err)
	}
	if _, err := w.Add(crashers); err != nil {
		return false, xerrors.Errorf("unable to add files: %w", err)
	}
	// Dirty check whether files were added.
	if s, err = w.Status(); err != nil || isClean(s) != nil {
		return true, nil
	}
	return false, nil
}

func isClean(s git.Status) error {
	for file, status := range s {
		if status.Staging != git.Unmodified && status.Staging != git.Untracked {
			return xerrors.Errorf("%c %s", status.Staging, file)
		}
	}
	return nil
}
