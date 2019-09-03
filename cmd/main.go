package main

import (
	"log"

	"github.com/Oncilla/fuzzinator/conf"
	"github.com/Oncilla/fuzzinator/execute"
)

func main() {
	commit, err := execute.CommitHash()
	if err != nil {
		log.Fatalf("Unable to get hash: %s", err)
	}
	target := conf.Target{
		Corpus: "./cmd/testdata/corpus",
		Name:   "fuzz",
		Harness: conf.Harness{
			Function: "Fuzz",
			Package:  "github.com/Oncilla/fuzzinator/cmd/fuzz",
		},
	}
	tmpDir, err := execute.SetupTempDir(target.Name, commit)
	if err != nil {
		log.Fatalf("Unable to setup temp dir: %s", err)
	}
	if err := execute.SetupCorpus(target.Corpus, tmpDir); err != nil {
		log.Fatalf("Unable to setup corpus: %s", err)
	}
	bin, err := execute.BuildBinary(target, tmpDir)
	if err != nil {
		log.Fatalf("Unable to build fuzzing binart: %s", err)
	}
	log.Println("Success", bin)
}
