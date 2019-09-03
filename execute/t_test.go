package execute

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gopkg.in/src-d/go-git.v4"
)

func TestA(t *testing.T) {
	// Clones the given repository, creating the remote, the local branches
	// and fetching the objects, everything in memory:
	r, err := git.PlainOpen("..")
	require.NoError(t, err)

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	require.NoError(t, err)
	t.Log(ref.Hash())

	// ... retrieves the commit history
	//	since := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	//	until := time.Date(2019, 7, 30, 0, 0, 0, 0, time.UTC)
	//	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
	//	require.NoError(t, err)

	// ... just iterates over the commits, printing it
	//	err = cIter.ForEach(func(c *object.Commit) error {
	//		fmt.Println(c)
	//
	//		return nil
	//	})
}
