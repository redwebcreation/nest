package pkg

import (
	"gotest.tools/v3/assert"
	"io/fs"
	"os"
	"testing"
)

func TestEnsureDirExists(t *testing.T) {
	assert.Assert(t, dirNotExists("/tmp/not-exists"))
	path := ensureDirExists("/tmp/not-exists")
	assert.Equal(t, path, "/tmp/not-exists")
	assert.Assert(t, dirExists("/tmp/not-exists"))
}

func TestEnsureDirExists2(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll("/tmp/not-exists")
		if err != nil {
			t.Errorf("Failed to remove /tmp/not-exists: %v", err)
		}
	})

	if _, err := os.Stat("/tmp/not-exists/config.json"); err == nil {
		t.Errorf("/tmp/not-exists/config.json should not exist (err: %v)", err)
	}

	path := ensureDirExists("/tmp/not-exists/config.json")
	assert.Assert(t, dirExists("/tmp/not-exists"))
	assert.Equal(t, path, "/tmp/not-exists/config.json")

	_, err := os.Stat("/tmp/not-exists/config.json")
	assert.ErrorIs(t, err, fs.ErrNotExist)
}

func dirExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func dirNotExists(path string) bool {
	return !dirExists(path)
}
