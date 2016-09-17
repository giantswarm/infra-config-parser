package cli

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// chdirTmp creates a temporary directory and changes the working directory to it.
// Returned clean function removes the temporary directory and reverts working directory.
// Must not be used in Parallel tests.
func chdirTmp(t *testing.T) (clean func()) {
	tmpdir, err := ioutil.TempDir(".", "tmp-test-")
	if err != nil {
		t.Fatalf("unexpected TempDir error = %v", err)
	}
	defer func() {
		if t.Failed() {
			os.RemoveAll(tmpdir)
		}
	}()
	if err := os.Chdir(tmpdir); err != nil {
		t.Fatalf("unexpected Chdir error = %v", err)
	}
	return func() {
		os.Chdir("..")
		os.RemoveAll(tmpdir)
	}
}

// prepareDir creates the group directory and files described by files argument.
func prepareDir(t *testing.T, index int, group string, files []fileDesc) {
	// Create group directory.
	if err := os.MkdirAll(group, os.FileMode(0755)); err != nil && group != "" {
		t.Fatalf("#%d: unexpected Mkdir error = %v", index, err)
	}
	// Create directories and file for each file's path.
	for _, f := range files {
		if i := strings.LastIndex(f.Path, "/"); i > 0 {
			if err := os.MkdirAll(f.Path[:i], os.FileMode(0755)); err != nil && !os.IsExist(err) {
				t.Fatalf("#%d: unexpected MkdirAll error = %v", index, err)
			}
		}
		if err := ioutil.WriteFile(f.Path, []byte(f.Content), os.FileMode(0644)); err != nil {
			t.Fatalf("#%d: unexpected WriteFile error = %v", index, err)
		}
	}
}