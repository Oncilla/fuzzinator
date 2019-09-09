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
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"

	"github.com/oncilla/fuzzinator/conf"
	"github.com/oncilla/fuzzinator/lib"
)

var (
	confFile  string
	terminate <-chan struct{}
)

var rootCmd = &cobra.Command{
	Use:   "fuzzinator",
	Short: "fuzzinator streamlines go fuzzing",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target, commit, err := targetAndCommit(confFile, args[0])
		if err != nil {
			return err
		}
		if err := setup(target, commit, terminate); err != nil {
			return err
		}
		if err := fuzz(target.Name, commit, terminate); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	terminate = handleSigTerm()

	rootCmd.PersistentFlags().StringVarP(&confFile, "conf", "c", "fuzzinator.yml",
		"defines the config file path (default fuzzinator.yml)")
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(fuzzCmd)
	rootCmd.AddCommand(crashersCmd)
}

// Execute executes the comands.
func Execute() {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func targetAndCommit(confFile, targetName string) (conf.Target, string, error) {
	var cfg conf.Conf
	raw, err := ioutil.ReadFile(confFile)
	if err != nil {
		return conf.Target{}, "", xerrors.Errorf("unable to read config file at %s: %w",
			confFile, err)
	}
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return conf.Target{}, "", xerrors.Errorf("unable to parse config file at %s: %w",
			confFile, err)
	}
	target, ok := cfg.Targets[targetName]
	if !ok {
		return conf.Target{}, "", xerrors.Errorf("target not in config file at %s: %w",
			confFile, err)
	}
	dir, err := lib.PkgDir(target.Harness.Package)
	if err != nil {
		return conf.Target{}, "", xerrors.Errorf("error resolving package %q: %w",
			target.Harness.Package, err)
	}
	commit, err := lib.CommitHash(dir)
	if err != nil {
		return conf.Target{}, "", xerrors.Errorf("unable to get git commit id: %w", err)
	}
	return target, commit, nil
}

func handleSigTerm() <-chan struct{} {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	stop := make(chan struct{})
	go func() {
		<-c
		close(stop)
		for range c {
			log.Fatalln("Received further SIGTERM before finishing execution")
		}
	}()
	return stop
}
