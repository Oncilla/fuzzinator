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

package conf

import "golang.org/x/xerrors"

// Conf configures the fuzzinator. It is designed to be compatible with
// fuzzbuzz.io project yaml.
type Conf struct {
	// Targets contains all fuzzing targets.
	Targets TargetMap `yaml:"targets"`
}

// TargetMap contains all targets and ensures no two targets can share the same
// name.
type TargetMap map[string]Target

// UnmarshalYAML translates the targets list to a map
func (m *TargetMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var list []Target
	if err := unmarshal(&list); err != nil {
		return xerrors.Errorf("unmarshalling list: %w", err)
	}
	*m = make(TargetMap)
	for _, target := range list {
		if _, ok := (*m)[target.Name]; ok {
			return xerrors.Errorf("target already exists: %s", target.Name)
		}
		(*m)[target.Name] = target
	}
	return nil
}

// Target defines a single fuzzing target.
type Target struct {
	Name     string  `yaml:"name"`
	Corpus   string  `yaml:"corpus"`
	Crashers string  `yaml:"crashers"`
	Harness  Harness `yaml:"harness"`
}

// Harness defines the fuzzing harness.
type Harness struct {
	// BuildTags contains the optional build tags. The 'gofuzz' build tag will
	// be set by fuzzinator itself.
	BuildTags string `yaml:"build_tags"`
	// Function specifies the entry point for fuzzing.
	Function string `yaml:"function"`
	// Package specifies the package of the entry point.
	Package string `yaml:"package"`
}
