package git

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func Test_LsTree(t *testing.T) {
	tree := Tree{
		dirPath: "./testdata",
		path:    ".",
		branch:  "main",
	}

	err := tree.LsTree()
	assert.Equal(t, nil, err)
}
