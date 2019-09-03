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

package conf_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Oncilla/fuzzinator/conf"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v2"
)

func TestCompatible(t *testing.T) {
	raw, err := ioutil.ReadFile("testdata/fuzzbuzz.yml")
	require.NoError(t, err)
	var cfg conf.Conf
	err = yaml.Unmarshal(raw, &cfg)
	require.NoError(t, err)
	yamlTarget := conf.Target{
		Name:   "FromYAML",
		Corpus: "./corpus",
		Harness: conf.Harness{
			Function: "FromYAML",
			Package:  "github.com/fuzzbuzz/tutorial",
		},
	}
	assert.Equal(t, yamlTarget, cfg.Targets[yamlTarget.Name])
	jsonTarget := conf.Target{
		Name:   "FromJSON",
		Corpus: "./corpus",
		Harness: conf.Harness{
			Function: "FromJSON",
			Package:  "github.com/fuzzbuzz/tutorial",
		},
	}
	assert.Equal(t, jsonTarget, cfg.Targets[jsonTarget.Name])
}
